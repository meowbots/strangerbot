package main

import "log"

// EndConversationEvent contains the data of the conversation that has to be
// ended
type EndConversationEvent struct {
	ChatID int64
}

func endConversationWorker(jobs <-chan EndConversationEvent) {
	for e := range jobs {
		u, err := retrieveUser(e.ChatID)

		if err != nil {
			log.Printf("Could not retrieve user in worker %s", err)
			return
		}

		// Check if is valid
		if u.MatchChatID.Valid {
			db.Exec("UPDATE users SET match_chat_id = NULL, available = 0, previous_match = ? WHERE chat_id = ?", u.MatchChatID, u.ChatID)
			db.Exec("UPDATE users SET match_chat_id = NULL, available = 0, previous_match = ? WHERE chat_id = ?", u.ChatID, u.MatchChatID)

			telegram.SendMessage(u.MatchChatID.Int64, "Sadly, we’re ending the conversation…Although time spent together was short, we still hope you enjoyed conversing with one another here in the TaveRHn! Type /start to get matched with a new hero!", emptyOpts)
			telegram.SendMessage(u.MatchChatID.Int64, "Type /start to get matched with a new hero", emptyOpts)
		} else {
			db.Exec("UPDATE users SET available = 0 WHERE chat_id = ?", u.ChatID)
		}

		telegram.SendMessage(u.ChatID, "Sadly, we’re ending the conversation…Although time spent together was short, we still hope you enjoyed conversing with one another here in the TaveRHn! Type /start to get matched with a new hero!", emptyOpts)
		telegram.SendMessage(u.ChatID, "Type /start to get matched with a new hero!", emptyOpts)
	}
}
