package main

import "strings"

var techWords = `
go fmt nil map var int byte rune func chan err ctx api cli git bug fix log
test sync http json yaml port host ping tcp udp sql cron heap stack slice
defer panic select struct const range break build merge rebase commit push
pull fork clone deploy docker server client proxy cache queue mutex token
lambda vector tensor kernel socket buffer stream thread router schema
pointer channel package import return switch method string float error
module vendor binary compile linker runtime context timeout latency
goroutine interface waitgroup benchmark profiling allocator scheduler
middleware container kubernetes terraform microservice serializer
refactor deadlock checksum hashmap iterator generics closure recursion
algorithm database migration callback promise async await unicode
pipeline observability concurrency throughput bottleneck segfault
monorepo dependency changelog hotfix rollback staging production
`

var everydayWords = `
sun sky sea air ice fog dew mud sand rock leaf tree root seed bark moss
rain snow wind hail mist storm cloud thunder lightning rainbow drizzle
river ocean lake pond creek wave tide shore beach island valley canyon
hill cliff cave desert forest jungle meadow prairie glacier volcano
flower petal blossom garden orchard harvest sunrise sunset twilight
moon star comet planet galaxy orbit eclipse meteor universe horizon

cat dog cow pig hen owl bat bee ant fox elk yak emu ram eel cod
bear wolf deer frog duck goat lion seal swan crab moth wasp toad
horse tiger zebra panda koala camel sheep mouse otter eagle raven
shark whale squid snail spider rabbit turtle parrot falcon dolphin
penguin giraffe leopard cheetah gorilla hamster lobster octopus
butterfly crocodile kangaroo flamingo porcupine hedgehog jellyfish

tea jam pie egg ham nut fig oat rice corn bean soup milk cake salt
bread honey lemon mango peach grape melon olive onion pasta pizza
apple berry cherry banana carrot potato tomato pepper cookie muffin
waffle pancake noodle cheese butter yogurt cereal biscuit pretzel
avocado broccoli cucumber pumpkin popcorn sandwich chocolate pineapple
strawberry blueberry watermelon cinnamon spaghetti marshmallow

red blue gold pink gray teal jade ruby amber coral ivory azure olive
green black white brown silver violet orange purple yellow crimson
scarlet emerald indigo magenta turquoise lavender maroon

home room door roof wall bed desk lamp sofa rug cup mug pan pot jar
chair table clock mirror window pillow blanket candle basket bucket
kettle teapot spoon fork knife plate bottle drawer ladder hammer
umbrella backpack suitcase notebook pencil scissors envelope keyboard
bicycle scooter tractor airplane rocket balloon lantern compass

city town road park farm shop bank mall zoo gym bus car van ship
train bridge tunnel castle temple palace museum library airport
station harbor market bakery cinema stadium village kingdom

run hop sit fly dig sing swim walk jump kick skip spin dive ride
read write draw paint dance laugh smile dream sleep climb throw catch
build cook bake wash fold sweep plant water carry listen whisper
wander explore discover imagine remember believe celebrate

happy angry brave calm shy kind bold wise cool warm cozy fuzzy
silly sleepy hungry thirsty curious gentle clever honest humble
cheerful grateful peaceful playful hopeful joyful fearless

big small tiny huge tall short wide deep high low fast slow loud
quiet soft hard dark light bright smooth rough heavy fluffy shiny
sparkly ancient modern golden frozen wild gentle mighty

love hope luck time life game song tune book poem tale myth idea
heart smile magic wonder secret puzzle riddle legend journey memory
friend family winter spring summer autumn holiday birthday weekend
morning evening midnight adventure treasure mystery freedom courage

guitar piano violin drum flute trumpet melody rhythm harmony chorus
soccer tennis hockey cricket chess karate surfing skating marathon

quiz jazz fizz buzz zigzag jigsaw jukebox oxygen zephyr quartz
whimsical serendipity kaleidoscope labyrinth extraordinary
spectacular wonderland breakfast fireworks lighthouse snowflake
thunderstorm watermelon dragonfly afternoon everywhere
`

var wordList = func() []string {
	seen := map[string]bool{}
	var out []string
	for _, w := range strings.Fields(techWords + everydayWords) {
		if !seen[w] {
			seen[w] = true
			out = append(out, w)
		}
	}
	return out
}()

func wordsWithLen(min, max int) []string {
	var out []string
	for _, w := range wordList {
		if len(w) >= min && len(w) <= max {
			out = append(out, w)
		}
	}
	return out
}
