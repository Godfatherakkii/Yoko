package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	tb "gopkg.in/tucnak/telebot.v2"
)

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}


func get_user(m *tb.Message) (string, string) {
	if m.IsReply() {
		user_obj, _ := strconv.Itoa(m.ReplyTo.Sender.ID)
		if len(m.Payload) != 0 {
			return user_obj, m.Payload
		} else {
			return user_obj, ""
		}
	} else if len(m.Payload) != 0 {
		x := strings.SplitN(m.Payload, " ", 2)
                if len(x) > 1 {
				return x[0], x[1]
                } else {
				return x[0], ""
		}
	} else {
                b.Reply(m, "You dont seem to be referring to a user or the ID specified is incorrect..")
		return "", ""
	}
}

func get_entity(m *tb.Message, user_id string) *tb.Chat {
 entity, err := b.ChatByID(user_id)
 if err != nil{
          b.Reply(m, "Looks like I don't have control over that user, or the ID isn't a valid one. If you reply to one of their messages, I'll be able to interact with them.")
          return nil
 }
 return entity
}
    

func info(m *tb.Message) {
	user_id, _ := get_user(m)
        if user_id == string(""){
            return b.Reply(m, "No user")
        }
        user_obj := get_entity(m, user_id)
	final_msg := fmt.Sprintf("<b>User info</b>\n<b>ID:</b> <code>%s</code>\n<b>First Name:</b> %s\n<b>Last Name:</b> %s\n<b>IsBot:</b> %s\n<b>Username:</b> @%s\n\n<b>Gbanned:</b> %s", strconv.Itoa(user_obj.ID), user_obj.FirstName, user_obj.LastName, strconv.FormatBool(user_obj.IsBot), user_obj.Username, "No")
	_, err := b.Reply(m, final_msg)
	if err != nil {
		fmt.Println(err)
	}
}
