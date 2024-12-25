package data

import (
	"github.com/moyai-network/teams/internal/model"
	"gopkg.in/square/go-jose.v2/json"
	"io"
	"net/http"
	"strconv"
)

var url = "https://minecraftpocket-servers.com/api/?object=servers&month=current&element=voters&key=0LnGcO8ioywfbZvrIoLM0abshjV4cdOmZuC&format=json"

type voterData struct {
	Nickname string `json:"nickname"`
	Votes    string `json:"votes"` // PocketServer is retarded and returns a string.
}

func (d voterData) VoteCount() int {
	n, _ := strconv.Atoi(d.Votes)
	return n
}

type voteData struct {
	Name    string      `json:"name"`
	Address string      `json:"address"`
	Port    string      `json:"port"`
	Month   string      `json:"month"`
	Voters  []voterData `json:"voters"`
}

func requestVoteData() (voteData, error) {
	var dat voteData
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dat, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return voteData{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dat, err
	}

	err = json.Unmarshal(body, &dat)
	return dat, err
}

func NewVoters() []model.User {
	/*dat, err := requestVoteData()
	if err != nil {
		return nil
	}

	var users []model.User
	for _, v := range dat.Voters {
		u, err := LoadUserFromName(v.Nickname)
		if err != nil || time.Since(u.LastVote) < 24*time.Hour {
			continue
		}
		u.LastVote = time.Now()
		users = append(users, u)
	}*/
	return []model.User{}
}
