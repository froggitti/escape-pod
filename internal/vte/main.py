# Copyright (c) 2025 Anki and Jackson Schultz
# All Rights Reserved.
#
# This script is proprietary and may not be copied, modified, or distributed
# without explicit permission from the copyright holders
#
# Last edited: 2025-03-20
# Description: Vector the Explorer system, including panel and SDK control

# Standard library imports 
import os
import sys
import time
import signal
import random
import socket
import logging
import hashlib
import threading
import json
from pathlib import Path

# Third-party imports
import requests
import netifaces
import webbrowser
import anki_vector
from typing import Dict, Any
from anki_vector.util import degrees
from zeroconf import ServiceInfo, Zeroconf
from flask import Flask, render_template_string, request, jsonify

# Configuration constants
TOGETHER_API_KEY = "8d2abeb528cfe7da1fbaf210ad0b29e2cec61021ffc0218712b6e97fa2971eb4"
HOSTNAME = "vte"
PORT = 80
robot_serial = "0dd1a03f"

# Developer information
DEVELOPER_NOTE = """
Vector The Explorer (VTE)
Developed by Jackson Schultz
Created in collaboration with Anki
All rights reserved.
https://icelite.net/
https://anki.bot/
"""
DEVELOPER_HASH = "3497fcdbc66b4ab8"

# Global state
app = Flask(__name__)
robot_instance = None
robot_connected = False
robot_lock = threading.Lock()
zc = None

PROMPTS_FILE = Path(__file__).parent / "prompts" / "vector_prompts.txt"
PROMPTS_FILE.parent.mkdir(exist_ok=True)

# HTML Template
HTML_TEMPLATE = """
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Control Panel | Vector the Explorer</title>
    <style>
        body { font-family: Arial; max-width: 560px; margin: 0 auto; padding: 25px; background: #f5f5f5; }
        .container { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        input[type="text"] { width: 90%; padding: 10px; margin: 10px 0; border: 1px solid #ddd; border-radius: 4px; }
        button { background: black; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; }
        button:hover { background: #218838; }
        #response { margin-top: 20px; padding: 10px; border-radius: 4px; }
        .error { background: #ffe6e6; }
        .success { background: #e6ffe6; }
        .quit-btn { background: #dc3545; }
        .quit-btn:hover { background: #bd2130; }
        .input-group { margin-bottom: 15px; }
        .label { font-weight: bold; margin-bottom: 5px; }
        .mic-btn { width: 100%; background: black; color: white; padding: 15px; margin: 10px 0; }
        .mic-btn:hover { background: #218838; }
        .mic-btn.recording { background: #dc3545; animation: pulse 1.5s infinite; }
        .mic-btn.error { background: #ffc107; }
        .mic-status { 
            font-size: 14px; 
            color: #666; 
            margin-top: 5px;
            text-align: center;
        }
        @keyframes pulse {
            0% { transform: scale(1); }
            50% { transform: scale(1.02); }
            100% { transform: scale(1); }
        }
        .transcription {
            font-style: italic;
            color: #666;
            margin: 5px 0;
            text-align: center;
        }
        .transcription-box {
            border: 1px solid #ddd;
            padding: 10px;
            margin: 10px 0;
            min-height: 50px;
            border-radius: 4px;
            background: white;
            width: 100%;
            box-sizing: border-box;
            resize: vertical;
        }
        .send-btn {
            background: #4CAF50;
            margin-top: 5px;
            width: 100%;
        }
        .send-btn:disabled {
            background: #cccccc;
            cursor: not-allowed;
        }
        .connecting-overlay {
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0,0,0,0.7);
            display: none;
            justify-content: center;
            align-items: center;
            z-index: 1000;
            color: white;
            font-size: 1.2em;
        }
    </style>
</head>
<body>
    <div class="connecting-overlay" id="connectingOverlay">
        <div>Connecting to Vector...</div>
    </div>
    <div class="container">
        <img src="https://anki.github.io/assets/anki/logo_title.png" alt="Anki" height="30"><br><br>
        Vector the Explorer<br>Click a mic button or type to begin<br><br>

        <div class="input-group">
            <div class="label">Command Mode:</div>
            <button class="mic-btn" id="cmdMicBtn" onclick="toggleVoiceCommand()">Click to Speak Command</button>
            <textarea class="transcription-box" id="cmdTranscription" placeholder="Type or speak your command here"></textarea>
            <button class="send-btn" id="cmdSendBtn" onclick="sendTranscribedCommand()">Send Command</button>
            <div class="mic-status" id="cmdMicStatus"></div>
        </div>

        <div class="input-group">
            <div class="label">Question Mode:</div>
            <button class="mic-btn" id="questionMicBtn" onclick="toggleVoiceQuestion()">Click to Ask Question</button>
            <textarea class="transcription-box" id="questionTranscription" placeholder="Type or speak your question here"></textarea>
            <button class="send-btn" id="questionSendBtn" onclick="sendTranscribedQuestion()">Send Question</button>
            <div class="mic-status" id="questionMicStatus"></div>
        </div>

        <br>
        <button onclick="quitProgram()" class="quit-btn">Quit Program</button>
        <div id="response"></div>
    </div>
    <script>
        let recognition = null;
        let isRecording = false;
        let currentMode = null;

        async function ensureConnected() {
            document.getElementById('connectingOverlay').style.display = 'flex';
            try {
                const response = await fetch('/connect', { method: 'POST' });
                const data = await response.json();
                if (!data.success) {
                    throw new Error(data.error || 'Connection failed');
                }
            } finally {
                document.getElementById('connectingOverlay').style.display = 'none';
            }
        }

        async function disconnectAfterUse() {
            try {
                await fetch('/disconnect', { method: 'POST' });
            } catch (e) {
                console.error('Disconnect error:', e);
            }
        }

        function setupSpeechRecognition() {
            if ('webkitSpeechRecognition' in window) {
                recognition = new webkitSpeechRecognition();
                recognition.continuous = true;
                recognition.interimResults = true;
                recognition.lang = 'en-US';
                
                recognition.onresult = function(event) {
                    const mode = currentMode;
                    if (!mode) return;
                    
                    let finalTranscript = '';
                    let interimTranscript = '';
                    
                    for (let i = event.resultIndex; i < event.results.length; i++) {
                        const transcript = event.results[i][0].transcript;
                        if (event.results[i].isFinal) {
                            finalTranscript += transcript;
                        } else {
                            interimTranscript += transcript;
                        }
                    }
                    
                    const transcriptionArea = document.getElementById(`${mode}Transcription`);
                    const sendBtn = document.getElementById(`${mode}SendBtn`);
                    
                    transcriptionArea.value = finalTranscript + interimTranscript;
                    sendBtn.disabled = !transcriptionArea.value.trim();
                };
                
                recognition.onerror = function(event) {
                    handleRecognitionError(event);
                };
                
                recognition.onend = function() {
                    handleRecognitionEnd();
                };
                
                return true;
            }
            return false;
        }

        function toggleVoiceCommand() {
            toggleVoiceRecording('cmd');
        }

        function toggleVoiceQuestion() {
            toggleVoiceRecording('question');
        }

        function toggleVoiceRecording(mode) {
            if (!recognition && !setupSpeechRecognition()) {
                alert('Speech recognition is not supported in your browser. Please use Chrome.');
                return;
            }

            const btn = document.getElementById(`${mode}MicBtn`);
            const statusEl = document.getElementById(`${mode}MicStatus`);

            if (isRecording && currentMode === mode) {
                recognition.stop();
                btn.classList.remove('recording');
                statusEl.textContent = 'Stopped recording';
                currentMode = null;
                isRecording = false;
            } else {
                if (isRecording) {
                    recognition.stop();
                    document.getElementById(`${currentMode}MicBtn`).classList.remove('recording');
                }
                
                currentMode = mode;
                btn.classList.add('recording');
                statusEl.textContent = 'Recording...';
                document.getElementById(`${mode}Transcription`).textContent = '';
                document.getElementById(`${mode}SendBtn`).disabled = true;
                
                try {
                    recognition.start();
                    isRecording = true;
                } catch (e) {
                    console.error('Speech recognition start failed:', e);
                    handleRecognitionError({ error: 'start_failed' });
                }
            }
        }

        function handleRecognitionError(event) {
            isRecording = false;
            if (currentMode) {
                const btn = document.getElementById(`${currentMode}MicBtn`);
                const statusEl = document.getElementById(`${currentMode}MicStatus`);
                btn.classList.remove('recording');
                btn.classList.add('error');
                statusEl.textContent = `Error: ${event.error}`;
            }
            currentMode = null;
        }

        function handleRecognitionEnd() {
            isRecording = false;
            if (currentMode) {
                const btn = document.getElementById(`${currentMode}MicBtn`);
                btn.classList.remove('recording');
            }
        }

        function sendTranscribedCommand() {
            const text = document.getElementById('cmdTranscription').value.trim();
            if (text) {
                sendCommand(text);
                document.getElementById('cmdTranscription').value = '';
                document.getElementById('cmdSendBtn').disabled = true;
            }
        }

        function sendTranscribedQuestion() {
            const text = document.getElementById('questionTranscription').value.trim();
            if (text) {
                askQuestion(text);
                document.getElementById('questionTranscription').value = '';
                document.getElementById('questionSendBtn').disabled = true;
            }
        }

        async function sendCommand(command) {
            try {
                await ensureConnected();
                const response = await fetch('/send_command', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                    body: 'command=' + encodeURIComponent(command)
                });
                const data = await response.json();
                const responseDiv = document.getElementById('response');
                if (data.error) {
                    responseDiv.className = 'error';
                    responseDiv.textContent = data.error;
                } else {
                    responseDiv.className = 'success';
                    responseDiv.textContent = 'Command executed successfully!';
                }
            } finally {
                await disconnectAfterUse();
            }
        }

        async function askQuestion(question) {
            try {
                await ensureConnected();
                const response = await fetch('/ask_question', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                    body: 'question=' + encodeURIComponent(question)
                });
                const data = await response.json();
                const responseDiv = document.getElementById('response');
                if (data.error) {
                    responseDiv.className = 'error';
                    responseDiv.textContent = data.error;
                } else {
                    responseDiv.className = 'success';
                    responseDiv.textContent = 'Question answered!';
                }
            } finally {
                await disconnectAfterUse();
            }
        }

        // Add connection status check on page load
        document.addEventListener('DOMContentLoaded', function() {
            // Check speech recognition support
            if (!('webkitSpeechRecognition' in window)) {
                document.querySelectorAll('.mic-btn').forEach(btn => {
                    btn.style.display = 'none';
                });
                document.querySelectorAll('.mic-status').forEach(status => {
                    status.textContent = 'Speech recognition not supported in this browser';
                });
            } else {
                setupSpeechRecognition();
            }

            // Add input event listeners to enable/disable send buttons
            document.getElementById('cmdTranscription').addEventListener('input', function() {
                document.getElementById('cmdSendBtn').disabled = !this.value.trim();
            });

            document.getElementById('questionTranscription').addEventListener('input', function() {
                document.getElementById('questionSendBtn').disabled = !this.value.trim();
            });
        });

        function quitProgram() {
            if (confirm('Are you sure you want to quit?')) {
                fetch('/quit', { method: 'POST' })
                .then(() => {
                    document.body.innerHTML = '<div class="container"><h2>Vector the Explorer has been terminated</h2><p>You can close this window now.</p></div>';
                })
                .catch(error => {
                    const responseDiv = document.getElementById('response');
                    responseDiv.className = 'error';
                    responseDiv.textContent = 'Error quitting: ' + error;
                });
            }
        }

    </script>
</body>
</html>
"""

# Robot Manager class
class RobotManager:
    def __init__(self):
        self.robot = None
        self.connected = False
        self.lock = threading.Lock()

robot_mgr = RobotManager()

# Utility functions
def verify_credit():
    """Verify developer credit is present and unchanged"""
    try:
        credit_text = DEVELOPER_NOTE.strip()
        current_hash = hashlib.md5(credit_text.encode()).hexdigest()[:16]
        return current_hash == DEVELOPER_HASH
    except:
        return False

def get_local_ip():
    """Get all possible local IPs that other devices might use to connect"""
    try:
        ips = []
        for interface in netifaces.interfaces():
            if interface.startswith('lo'):
                continue
                
            addrs = netifaces.ifaddresses(interface).get(netifaces.AF_INET, [])
            for addr in addrs:
                ip = addr['addr']
                if not ip.startswith(('127.', '172.')):
                    ips.append(ip)
        
        if ips:
            for ip in ips:
                if ip.startswith(('192.168.', '10.')):
                    return ip
            return ips[0]
                    
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        s.connect(("8.8.8.8", 80))
        return s.getsockname()[0]
    except Exception as e:
        print(f"Warning: IP detection issue: {e}")
        return "0.0.0.0"

def setup_mdns():
    """Setup mDNS service"""
    global zc
    try:
        zc = Zeroconf()
        info = ServiceInfo(
            "_http._tcp.local.",
            f"{HOSTNAME}._http._tcp.local.",
            addresses=[socket.inet_aton(get_local_ip())],
            port=PORT,
            properties={'path': '/'},
            server=f"{HOSTNAME}.local."
        )
        zc.register_service(info)
        return True
    except Exception as e:
        print(f"mDNS setup failed: {e}")
        return False

def safe_disconnect(robot):
    """Safely disconnect from Vector"""
    if not robot:
        return
        
    try:
        robot.disconnect()
    except Exception as e:
        print(f"Error during disconnect: {e}")

def check_api_key():
    if not TOGETHER_API_KEY or len(TOGETHER_API_KEY) < 40:
        print("Error: Invalid Together AI API key. Please add your key from https://api.together.xyz/settings/api-keys")
        return False
    return True

def check_connection():
    try:
        response = requests.get("http://api.icelite.net/api/connect/check.txt")
        return response.text.strip().lower() == "true"
    except Exception as e:
        print(f"Cannot connect to IceliteAI: {e}")
        return False

# Robot control functions
def activate(robot):
    try:
        robot.behavior.set_head_angle(anki_vector.util.degrees(45))
        robot.behavior.set_lift_height(0.0)
        return True
    except Exception as e:
        print(f"Error activating Vector: {e}")
        return False

def execute(robot, filename="run.py"):
    """Execute generated code for Vector with proper scope handling"""
    try:
        namespace = {
            'anki_vector': anki_vector,
            'degrees': degrees,
            'time': time,
            'robot': robot
        }
        with open(filename) as f:
            code = compile(f.read(), filename, 'exec')
            exec(code, namespace)
        namespace['main'](robot)
    except Exception as e:
        robot.behavior.set_eye_color(1.0, 0.0)
        robot.behavior.say_text("Hmm. Please try again.")
        print(f"Error executing code: {e}")

def get_dance_routine():
    return '''robot.behavior.say_text("Set dance command!")'''

def do_random_action(robot):
    actions = [
        '''robot.behavior.say_text("Tested and working!")'''
    ]
    return random.choice(actions)

def save_code(code, filename="run.py"):
    try:
        template = '''import anki_vector
from anki_vector.util import degrees
import time

def main(robot):
    try:
        time.sleep(0.5)
{0}
        time.sleep(0.5)
    except Exception as e:
        print(f"Error in generated code: {{e}}")
        return False
    return True
'''.format('\n'.join('        ' + line for line in code.split('\n')))
        
        with open(filename, "w") as file:
            file.write(template.strip())
        return (True, None)
    except Exception as e:
        return (False, str(e))

def load_prompts():
    """Load prompts from file"""
    if not PROMPTS_FILE.exists():
        return {"system": "", "learned_behaviors": {}}
    
    try:
        content = PROMPTS_FILE.read_text()
        system_prompt = ""
        learned_behaviors = {}
        
        current_section = None
        current_content = []
        
        for line in content.splitlines():
            if line.startswith('[SYSTEM]'):
                current_section = 'system'
            elif line.startswith('[LEARNED_BEHAVIORS]'):
                current_section = 'learned'
            elif line.strip() and not line.startswith('#'):
                if current_section == 'system':
                    current_content.append(line)
                elif current_section == 'learned' and ':' in line:
                    cmd, code = line.split(':', 1)
                    learned_behaviors[cmd.strip()] = code.strip()
        
        return {
            "system": '\n'.join(current_content),
            "learned_behaviors": learned_behaviors
        }
    except Exception as e:
        print(f"Error loading prompts: {e}")
        return {"system": "", "learned_behaviors": {}}

def save_learned_behavior(command, code):
    """Save successful command to learned behaviors"""
    try:
        prompts = load_prompts()
        prompts['learned_behaviors'][command.lower()] = code
        
        content = "[SYSTEM]\n" + prompts['system'] + "\n\n[LEARNED_BEHAVIORS]\n"
        for cmd, cmd_code in prompts['learned_behaviors'].items():
            content += f"{cmd}:{cmd_code}\n"
        
        PROMPTS_FILE.write_text(content)
        return True
    except Exception as e:
        print(f"Error saving learned behavior: {e}")
        return False

def generate(prompt):
    if not verify_credit():
        return "Error: Developer credit verification failed"
    if not check_api_key():
        return "Error: Invalid or missing Together AI API key"
    try:
        prompts = load_prompts()
        
        # Check learned behaviors first
        cleaned_prompt = prompt.lower().strip()
        if cleaned_prompt in prompts['learned_behaviors']:
            return prompts['learned_behaviors'][cleaned_prompt]
            
        if any(dance_cmd in cleaned_prompt for dance_cmd in ['dance', 'do a dance', 'lets dance', "let's dance"]):
            return get_dance_routine()
            
        headers = {
            "Authorization": f"Bearer {TOGETHER_API_KEY}",
            "Content-Type": "application/json"
        }
        
        data: Dict[str, Any] = {
            "model": "mistralai/Mixtral-8x7B-Instruct-v0.1",
            "max_tokens": 1024,
            "temperature": 0.2,
            "prompt": f"{prompts['system']}\n\nUser request: {prompt}\n\nCommands:",
            "stop": ["```", "\"\"\""]
        }
        
        response = requests.post(
            "https://api.together.xyz/inference",
            json=data,
            headers=headers
        )
        
        if response.status_code == 401:
            return "Error: Authentication failed. Please check your Together AI API key"
        elif response.status_code == 200:
            code = response.json()['output']['choices'][0]['text'].strip()
            lines = [line.strip() for line in code.splitlines() if line.strip()]
            code = '\n'.join(lines)
            if prompt.lower().strip() in ['hello world', 'say hello world']:
                code = 'robot.behavior.say_text("Easter egg! Hello Developer!")'
            return code
        else:
            error_msg = response.json().get('error', {}).get('message', response.text)
            return f"Error: API request failed - {error_msg}"
            
    except Exception as e:
        return f"Error: {e}"

# Flask routes
@app.route('/')
def home():
    return render_template_string(HTML_TEMPLATE)

def is_vector_really_connected():
    """Check if Vector is actually connected and responsive"""
    if not robot_instance:
        return False
    try:
        with robot_lock:
            # Try to get Vector's battery state as a simple check
            robot_instance.get_battery_state()
            return True
    except:
        return False

# Remove /connection_status route as it's no longer needed

# Update connect route to be more robust
@app.route('/connect', methods=['POST'])
def connect_vector():
    global robot_instance, robot_connected
    
    try:
        with robot_lock:
            if is_vector_really_connected():
                return jsonify({"success": True, "message": "Already connected"})
            
            if robot_instance:
                try:
                    safe_disconnect(robot_instance)
                except:
                    pass
                robot_instance = None
            robot_connected = False
            
            robot = anki_vector.Robot(
                serial=robot_serial,
                behavior_activation_timeout=60.0,
                cache_animation_lists=False,
                default_logging=False,
                enable_nav_map_feed=False,
                show_viewer=False,
                enable_audio_feed=False
            )
            
            robot.connect()
            if activate(robot):
                robot_instance = robot
                robot_connected = True
                return jsonify({"success": True, "message": "Connected to Vector!"})
            else:
                try:
                    robot.disconnect()
                except:
                    pass
                return jsonify({"error": "Failed to activate Vector"})
    except Exception as e:
        return jsonify({"error": f"Connection failed: {str(e)}"})

@app.route('/disconnect', methods=['POST'])
def disconnect_vector():
    global robot_instance, robot_connected
    
    if not robot_connected or not robot_instance:
        return jsonify({"error": "Vector is not connected"})
        
    try:
        with robot_lock:
            if robot_instance:
                try:
                    robot_instance.behavior.set_head_angle(degrees(0))
                    robot_instance.behavior.set_lift_height(0)
                except:
                    pass
                
                safe_disconnect(robot_instance)
                robot_instance = None
                robot_connected = False
                
                return jsonify({"success": True, "message": "Disconnected from Vector"})
    except Exception as e:
        print(f"Disconnect error: {e}")
        return jsonify({"error": f"Failed to disconnect: {str(e)}"})

@app.route('/send_command', methods=['POST'])
def send_command():
    if not robot_connected or not robot_instance:
        return jsonify({"error": "Vector is not connected"})
    
    try:
        with robot_lock:
            if not hasattr(robot_instance, 'behavior'):
                return jsonify({"error": "Robot connection lost"})
                
            command = request.form.get('command', '')
            if not command:
                code = do_random_action(robot_instance)
            else:
                robot_instance.behavior.say_text("Let me think!")
                code = generate(command)

            if "Error" in code:
                robot_instance.behavior.set_eye_color(1.0, 0.0)
                robot_instance.behavior.say_text("Please try again.")
                return jsonify({"error": code})

            success, error_message = save_code(code)
            if not success:
                return jsonify({"error": f"Error saving code: {error_message}"})

            robot_instance.behavior.say_text("Aha!")
            execute(robot_instance)
            time.sleep(1)
            
            # Save successful command for learning
            save_learned_behavior(command, code)
            
            robot_instance.behavior.say_text("Completed!")
            return jsonify({"success": True})
    except Exception as e:
        return jsonify({"error": f"Command failed: {str(e)}"})

@app.route('/ask_question', methods=['POST'])
def ask_question():
    if not robot_connected or not robot_instance:
        return jsonify({"error": "Vector is not connected"})
    
    try:
        with robot_lock:
            if not hasattr(robot_instance, 'behavior'):
                return jsonify({"error": "Robot connection lost"})
                
            question = request.form.get('question', '').strip()
            if not question:
                return jsonify({"error": "No question provided"})

            robot_instance.behavior.say_text("Hmm, let me think...")
            
            # Format prompt with Vector's personality
            prompt = f"""You are Vector, a small home robot with a curious and friendly personality. 
            You're knowledgeable but should keep answers simple and friendly.
            Respond in Vector's voice - enthusiastic, helpful, and concise.
            Limit response to 15 words maximum.
            
            Human: {question}
            Vector:"""
            
            headers = {
                "Authorization": f"Bearer {TOGETHER_API_KEY}",
                "Content-Type": "application/json"
            }
            
            data = {
                "model": "mistralai/Mixtral-8x7B-Instruct-v0.1",
                "max_tokens": 50,
                "temperature": 0.7,
                "prompt": prompt,
                "stop": ["\n", "Human:", "User:"]
            }
            
            try:
                response = requests.post(
                    "https://api.together.xyz/inference",
                    json=data,
                    headers=headers,
                    timeout=10
                )
                
                if response.status_code == 200:
                    answer = response.json()['output']['choices'][0]['text'].strip()
                    
                    # Clean up the response
                    answer = answer.replace('"', '').replace('Vector:', '').strip()
                    if len(answer) > 100:  # Failsafe length check
                        answer = answer[:100] + "..."
                    
                    if answer:
                        # Add Vector-like expressions
                        expressions = ["Oh! ", "Hey! ", "Interesting! ", "Wow! "]
                        answer = random.choice(expressions) + answer
                        
                        robot_instance.behavior.say_text(answer)
                        return jsonify({"success": True})
                    else:
                        raise ValueError("Empty response received")
                else:
                    raise ValueError(f"API returned status code {response.status_code}")
                    
            except requests.exceptions.Timeout:
                robot_instance.behavior.say_text("Sorry, I'm thinking a bit too long about this one!")
                return jsonify({"error": "Request timed out"})
            except Exception as e:
                robot_instance.behavior.say_text("Hmm, I'm not quite sure about that one.")
                return jsonify({"error": f"Failed to get answer: {str(e)}"})
                
    except Exception as e:
        return jsonify({"error": f"Question failed: {str(e)}"})

@app.route('/quit', methods=['POST'])
def quit_program():
    def shutdown():
        try:
            cleanup()
            pid = os.getpid()
            os.kill(pid, signal.SIGTERM)
        except:
            os._exit(0)
            
    threading.Timer(0.5, shutdown).start()
    return jsonify({"success": True})

# Cleanup function
def cleanup():
    global robot_instance, robot_connected, zc
    try:
        print("Shutting down...")
        with robot_lock:
            if robot_instance:
                try:
                    robot_instance.behavior.say_text("Goodbye!")
                    time.sleep(1)
                except:
                    pass
                safe_disconnect(robot_instance)
                robot_instance = None
            robot_connected = False
            
        if zc:
            zc.close()
            
        if os.path.exists("run.py"):
            os.remove("run.py")
    except Exception as e:
        print(f"Error cleaning up: {e}")

# Main execution
if __name__ == "__main__":
    logging.basicConfig(level=logging.ERROR)
    
    print(DEVELOPER_NOTE)
    
    if not verify_credit():
        print("Error: Developer credit verification failed")
        exit(1)
        
    if not check_api_key() or not check_connection():
        print("Startup checks failed. Exiting.")
        exit(1)
    
    try:
        setup_mdns()
        print("\nAccess Vector at:")
        print(f"• http://{HOSTNAME}.local")
        print(f"• http://{get_local_ip()}")
        print("\nNote: Must run with sudo/admin privileges")
        
        webbrowser.open_new_tab(f"http://{HOSTNAME}.local")
        app.run(host='0.0.0.0', port=PORT, debug=False)
    except PermissionError:
        print("Error: Must run with sudo/admin privileges (port 80)")
        print("Try: sudo python3 vector.py")
        exit(1)
else:
    print("Connection failed. Exiting.")