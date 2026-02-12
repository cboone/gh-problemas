package data

// UserClient resolves user-related data via GraphQL.
type UserClient struct {
	querier Querier
}

// NewUserClient creates a UserClient.
func NewUserClient(q Querier) *UserClient {
	return &UserClient{querier: q}
}

// WhoAmI returns the authenticated user's login.
func (c *UserClient) WhoAmI() (string, error) {
	var resp struct {
		Viewer struct {
			Login string `json:"login"`
		} `json:"viewer"`
	}

	if err := c.querier.Do("query { viewer { login } }", nil, &resp); err != nil {
		return "", err
	}

	return resp.Viewer.Login, nil
}
