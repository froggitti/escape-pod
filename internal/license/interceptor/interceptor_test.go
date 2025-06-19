package interceptor

func newInterceptor(lm LicenseManager) (*Interceptor, error) {
	return New("test_key_for_testing", lm)
}
