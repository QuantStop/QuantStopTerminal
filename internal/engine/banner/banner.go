package banner

import (
	"math/rand"
	"time"
)

var Banners []string

const (
	Graffiti = `
________                       __            __              ___________                  .__              .__   
\_____  \  __ _______    _____/  |_  _______/  |_  ____ _____\__    ___/__________  _____ |__| ____ _____  |  |  
 /  / \  \|  |  \__  \  /    \   __\/  ___/\   __\/  _ \\____ \|    |_/ __ \_  __ \/     \|  |/    \\__  \ |  |  
/   \_/.  \  |  // __ \|   |  \  |  \___ \  |  | (  <_> )  |_> >    |\  ___/|  | \/  Y Y  \  |   |  \/ __ \|  |__
\_____\ \_/____/(____  /___|  /__| /____  > |__|  \____/|   __/|____| \___  >__|  |__|_|  /__|___|  (____  /____/
       \__>          \/     \/          \/              |__|              \/            \/        \/     \/      
`
	BigMoney = `
  /$$$$$$                              /$$              /$$                   /$$$$$$$$                            /$$                  /$$
 /$$__  $$                            | $$             | $$                  |__  $$__/                           |__/                 | $$
| $$  \ $$/$$   /$$ /$$$$$$ /$$$$$$$ /$$$$$$  /$$$$$$$/$$$$$$   /$$$$$$  /$$$$$$| $$ /$$$$$$  /$$$$$$ /$$$$$$/$$$$ /$$/$$$$$$$  /$$$$$$| $$
| $$  | $| $$  | $$|____  $| $$__  $|_  $$_/ /$$_____|_  $$_/  /$$__  $$/$$__  $| $$/$$__  $$/$$__  $| $$_  $$_  $| $| $$__  $$|____  $| $$
| $$  | $| $$  | $$ /$$$$$$| $$  \ $$ | $$  |  $$$$$$  | $$   | $$  \ $| $$  \ $| $| $$$$$$$| $$  \__| $$ \ $$ \ $| $| $$  \ $$ /$$$$$$| $$
| $$/$$ $| $$  | $$/$$__  $| $$  | $$ | $$ /$\____  $$ | $$ /$| $$  | $| $$  | $| $| $$_____| $$     | $$ | $$ | $| $| $$  | $$/$$__  $| $$
|  $$$$$$|  $$$$$$|  $$$$$$| $$  | $$ |  $$$$/$$$$$$$/ |  $$$$|  $$$$$$| $$$$$$$| $|  $$$$$$| $$     | $$ | $$ | $| $| $$  | $|  $$$$$$| $$
 \____ $$$\______/ \_______|__/  |__/  \___/|_______/   \___/  \______/| $$____/|__/\_______|__/     |__/ |__/ |__|__|__/  |__/\_______|__/
      \__/                                                             | $$                                                                
                                                                       | $$                                                                
                                                                       |__/                                                                
`
	Cards = `
.------.------.------.------.------.------.------.------.------.------.------.------.------.------.------.------.------.
|Q.--. |U.--. |A.--. |N.--. |T.--. |S.--. |T.--. |O.--. |P.--. |T.--. |E.--. |R.--. |M.--. |I.--. |N.--. |A.--. |L.--. |
| (\/) | (\/) | (\/) | :(): | :/\: | :/\: | :/\: | :/\: | :/\: | :/\: | (\/) | :(): | (\/) | (\/) | :(): | (\/) | :/\: |
| :\/: | :\/: | :\/: | ()() | (__) | :\/: | (__) | :\/: | (__) | (__) | :\/: | ()() | :\/: | :\/: | ()() | :\/: | (__) |
| '--'Q| '--'U| '--'A| '--'N| '--'T| '--'S| '--'T| '--'O| '--'P| '--'T| '--'E| '--'R| '--'M| '--'I| '--'N| '--'A| '--'L|
'------'------'------'------'------'------'------'------'------'------'------'------'------'------'------'------'------'
`
	Crawford2 = `
  ___  __ __  ____ ____  ______  ___________  ___  ____  ______   ___ ____  ___ ___ ____ ____   ____ _     
 /   \|  |  |/    |    \|      |/ ___|      |/   \|    \|      | /  _|    \|   |   |    |    \ /    | |    
|     |  |  |  o  |  _  |      (   \_|      |     |  o  |      |/  [_|  D  | _   _ ||  ||  _  |  o  | |    
|  Q  |  |  |     |  |  |_|  |_|\__  |_|  |_|  O  |   _/|_|  |_|    _|    /|  \_/  ||  ||  |  |     | |___ 
|     |  :  |  _  |  |  | |  |  /  \ | |  | |     |  |    |  | |   [_|    \|   |   ||  ||  |  |  _  |     |
|     |     |  |  |  |  | |  |  \    | |  | |     |  |    |  | |     |  .  |   |   ||  ||  |  |  |  |     |
 \__,_|\__,_|__|__|__|__| |__|   \___| |__|  \___/|__|    |__| |_____|__|\_|___|___|____|__|__|__|__|_____|
`
	Graceful = `
  __  _  _  __  __ _ ____ ____ ____ __ ____ ____ ____ ____ _  _ __ __ _  __  __   
 /  \/ )( \/ _\(  ( (_  _/ ___(_  _/  (  _ (_  _(  __(  _ ( \/ (  (  ( \/ _\(  )  
(  O ) \/ /    /    / )( \___ \ )((  O ) __/ )(  ) _) )   / \/ \)(/    /    / (_/\
 \__\\____\_/\_\_)__)(__)(____/(__)\__(__)  (__)(____(__\_\_)(_(__\_)__\_/\_\____/
`
	Modular = `
 _______ __   __ _______ __    _ _______ _______ _______ _______ _______ _______ _______ ______   __   __ ___ __    _ _______ ___     
|       |  | |  |   _   |  |  | |       |       |       |       |       |       |       |    _ | |  |_|  |   |  |  | |   _   |   |    
|   _   |  | |  |  |_|  |   |_| |_     _|  _____|_     _|   _   |    _  |_     _|    ___|   | || |       |   |   |_| |  |_|  |   |    
|  | |  |  |_|  |       |       | |   | | |_____  |   | |  | |  |   |_| | |   | |   |___|   |_||_|       |   |       |       |   |    
|  |_|  |       |       |  _    | |   | |_____  | |   | |  |_|  |    ___| |   | |    ___|    __  |       |   |  _    |       |   |___ 
|      ||       |   _   | | |   | |   |  _____| | |   | |       |   |     |   | |   |___|   |  | | ||_|| |   | | |   |   _   |       |
|____||_|_______|__| |__|_|  |__| |___| |_______| |___| |_______|___|     |___| |_______|___|  |_|_|   |_|___|_|  |__|__| |__|_______|
`
	Soft = `
 .-----.                         .--.        .--.            .--------.                   .--.             .--. 
'  .-.  ' .--..--..--.--.--.--..-'  '-..---.-'  '-..---. .---'--.  .--,---.,--.--,--,--,--'--'--,--, ,--,--|  |
|  | |  | |  ||  ' .-.  |      '-.  .-(  .-'-.  .-| .-. | .-. | |  | | .-. |  .--|        |  |      ' .-.  |  |
'  '-'  '-'  ''  \ '-'  |  ||  | |  | .-'  ')|  | ' '-' | '-' ' |  | \   --|  |  |  |  |  |  |  ||  \ '-'  |  | 
 '-----'--''----' '--'--'--''--' '--' '----' '--'  '---'|  |-'  '--'  '----'--'  '--'--'--'--'--''--''--'--'--' 
                                                        '--'                                                    
`
	Varsity = `
                                _        _               _________                         _                __   
 . --- .                       / |_     / |_            |  _   _  |                       (_)              [  |
/  .-.  \ __   _  .--.  _ .--.'| |-'--.'| |-'.--. _ .--|_/ | | \_.---. _ .--. _ .--..--.  __  _ .--.  ,--.  | |  
| |   | |[  | | | '_\ : | .-. || |( ('\]| |/ .' \[ '/' \ \ | |  / /__\[ '/''\[ '.-. .-. |[  |[  .-. | '_\ : | |
\  '-'  \_| \_/ |// | |,| | | || |,''.'.| || \__.|| \__/ |_| |_ | \__.,| |    | | | | | | | | | | | |// | |,| |
'.___.\__'.__.'_\'-;__[___||__\__[\__) \__/'.__.'| ;.__/|_____| '.__.[___]  [___||__||__[___[___||__\'-;__[___]
                                                 |__|
`
)

func init() {
	Banners = []string{
		Graffiti,
		BigMoney,
		Cards,
		Crawford2,
		Graceful,
		Modular,
		Soft,
		Varsity,
	}
}

func GetRandomBanner() string {
	if len(Banners) <= 0 {
		return Graffiti
	}
	rand.Seed(time.Now().UnixNano())
	randBan := rand.Intn(len(Banners)-0) + 0
	return Banners[randBan]
}
