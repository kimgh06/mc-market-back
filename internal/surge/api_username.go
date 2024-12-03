package surge

import (
	"encoding/json"
	"fmt"
	"io"
	"maple/internal/utilities"
	"strconv"
	"strings"
)

func (a *API) ResolveUsernames(users []uint64) ([]string, error) {
	if len(users) < 1 {
		return []string{}, nil
	}

	formattedIdList := make([]string, len(users))
	for i := range users {
		formattedIdList[i] = strconv.FormatUint(users[i], 10)
	}

	joinedIdList := strings.Join(formattedIdList, ",")

	url := fmt.Sprintf("%s/v1/usernames?id_list=%s", a.Configuration.SurgeURL, joinedIdList)

	response, err := a.Client.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("%s", response.Status)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	var decoded *[]string
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, err
	}

	if decoded == nil {
		return []string{}, nil
	}

	return *decoded, nil
}

func (a *API) ResolveUsernamesAsMap(users []uint64) (map[uint64]string, error) {
	uniqueUsers := utilities.RemoveDuplicate(users)
	usernames, err := a.ResolveUsernames(uniqueUsers)
	if err != nil {
		return nil, err
	}

	result := make(map[uint64]string)

	for i := range uniqueUsers {
		result[uniqueUsers[i]] = usernames[i]
	}

	return result, nil
}
