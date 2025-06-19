let rsconf = {
  _id: "rs0",
  members: [
    {
     _id: 0,
     host: "localhost:27017",
     priority: 1,
    },
   ]
};
rs.initiate(rsconf);