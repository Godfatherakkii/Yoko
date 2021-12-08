package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
        "net/http"
        "time"
        "encoding/json"
        "github.com/StalkR/imdb"
	tb "gopkg.in/tucnak/telebot.v2"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

type MovieClient struct {
	Client *http.Client
}

func GetNewClient() *MovieClient {
	var mc *MovieClient
	mc = new(MovieClient)
	mc.Client = &http.Client{}
	return mc
}

func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
    if condition {
        return a
    }
    return b
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func get_user(m *tb.Message) (*tb.User, string) {
	if m.IsReply() {
		user_obj := m.ReplyTo.Sender
		if len(m.Payload) != 0 {
			return user_obj, m.Payload
		} else {
			return user_obj, ""
		}
	} else if len(m.Payload) != 0 {
		x := strings.SplitN(m.Payload, " ", 2)
		if isInt(x[0]) {
			user_obj, err := b.ChatByID(x[0])
                        if err != nil{
                                b.Reply(m, "Looks like I don't have control over that user, or the ID isn't a valid one. If you reply to one of their messages, I'll be able to interact with them.")
                                return nil, ""
                        }
                        user := &tb.User{ID: int(user_obj.ID), FirstName: user_obj.FirstName, LastName: user_obj.LastName, Username: user_obj.Username}
			if len(x) > 1 {
				return user, x[1]
			} else {
				return user, ""
			}
		} else {
                        u, err := getJson(m.Payload)
                        if err != nil {
                           b.Reply(m, fmt.Sprint(err.Error()))
                           return nil, ""
                        }
			user_obj := &tb.User{ID: int(u["id"].(float64)), Username: u["username"].(string), FirstName: u["first_name"].(string), LastName: u["last_name"].(string)}
			if len(x) > 1 {
				return user_obj, x[1]
			} else {
				return user_obj, ""
			}
		}
	} else {
                b.Reply(m, "You dont seem to be referring to a user or the ID specified is incorrect..")
		return nil, ""
	}
}

func get_perm(chat *tb.Chat, user *tb.User) {
 b.ChatMemberOf(chat, user)
}

func get_entity(m *tb.Message, user_id string) *tb.Chat {
 entity, err := b.ChatByID(user_id)
 if err != nil{
          b.Reply(m, "Looks like I don't have control over that user, or the ID isn't a valid one. If you reply to one of their messages, I'll be able to interact with them.")
          return nil
 }
 return entity
}

type mapType map[string]interface{}

func (t mapType) m(s string) mapType {
   return t[s].(map[string]interface{})
}

func getJson(url string) (mapType, error) {
    resp, err := myClient.Get("https://roseflask.herokuapp.com/username?username=" + url)
    if err != nil {
        fmt.Println("No response from request")
        return nil, err
    }
    defer resp.Body.Close()
    var t mapType
    json.NewDecoder(resp.Body).Decode(&t)   
    return t, err
}

func info(m *tb.Message) {
        if !m.IsReply() && string(m.Payload) == string(""){
            user_obj := m.Sender
            final_msg := fmt.Sprintf("<b>User info</b>\n<b>ID:</b> <code>%s</code>\n<b>First Name:</b> %s\n<b>Last Name:</b> %s\n<b>Username:</b> @%s\n<b>User Link:</b> <a href='tg://user?id=%s'>%s</a>\n\n<b>Gbanned:</b> %s", strconv.Itoa(int(user_obj.ID)), user_obj.FirstName, user_obj.LastName, user_obj.Username, strconv.Itoa(int(user_obj.ID)), "link", "No")
	    b.Reply(m, final_msg)
        } else {
          user_obj, _ := get_user(m)
          final_msg := fmt.Sprintf("<b>User info</b>\n<b>ID:</b> <code>%s</code>\n<b>First Name:</b> %s\n<b>Last Name:</b> %s\n<b>Username:</b> @%s\n<b>User Link:</b> <a href='tg://user?id=%s'>%s</a>\n\n<b>Gbanned:</b> %s", strconv.Itoa(int(user_obj.ID)), user_obj.FirstName, user_obj.LastName, user_obj.Username, strconv.Itoa(int(user_obj.ID)), "link", "No")
      	  b.Reply(m, final_msg)
        }
}

func gp(m *tb.Message) {
 u, _ := get_user(m)
 x, err := b.ChatMemberOf(m.Chat, u)
 fmt.Println(x.Rights)
 if err != nil {
    b.Reply(m, string(err.Error()))
    return 
 }
 b.Reply(m, fmt.Sprint(x.Rights))
}

type MovieInfo struct {
	Name        string
	Rating      string
	Description string
	PosterLink  string
}

func IMDb(m *tb.Message) {
 client := http.DefaultClient
 results, _ := imdb.SearchTitle(client, m.Payload)
 title, _ := imdb.NewTitle(client, results[0].ID)
 movie := fmt.Sprintf("<b>%s</b>\n<b>Type:</b> %s\n<b>Year:</b> %s", title.Name, title.Type, title.Year)
 b.Reply(m, movie)
}
