package main

import (
	"math/rand"
)

var killStatements = []string{
	" teabagged a piranha tank",
	" died of a heart attack while watching hentai",
	" just got back from yiffing",
	" was bitten by a horse",
	" was bent over and given a slap on the bottom",
	" changed their major to BIS",
	" got drilled",
	" got paddeled",
	" drank bleach",
	" was put on a group project with kleiderar",
	" was forever shunned by the resume gods",
	" became PayPay's bitch",
	" was yiffed by Chan",
	" was ignored by Eli",
	" was kicked from resistance",
	" had their appendix try to kill them",
	" voiced their political, religious, or other personal beliefs in resistance",
	" was punished for not conforming",
	" is the filling to a Jane Hansen sandwich",
	" suffocated in an amazon prime package",
	" had their privates waxed with ducktape",
	" was sent to Division 1 ICPC",
	" took ALL of Ritters Monolithic Kernel",
	" got unzipped by The Rit",
	" was strictly spoken to",
	" cut themselves on Chandler's edge",
	" was sacrificed to the void",
	" got baited haha XD",
	" got /pushups'd on",
	" opened Resistance after 5 PM",
	" choked on a Tide Pod",
	" is handing out with Logan Paul",
}

var killCommand = BotCommand{
	Name: "kill",
	Description: "insult someone, /kill my 8am",
	Matcher: messageContainsCommandMatcher("kill"),
	Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		commands := getContentFromCommand(update.Message.Text, "kill")
		if commands != "" {
			respChan <- *NewTextBotResponse(commands+killStatements[rand.Int() % len(killStatements)], update.Message.Chat.ID)
		}
	},
}