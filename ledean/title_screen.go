package ledean

import "fmt"

//GetStartScreen - Ascii art Start screen info as string
func GetStartScreen() string {
	//http://patorjk.com/software/taag/#p=display&f=Big%20Money-ne&t=LEDean
	return fmt.Sprintf(`

   /$$       /$$$$$$$$ /$$$$$$$
   | $$      | $$_____/| $$__  $$
   | $$      | $$      | $$  \ $$  /$$$$$$   /$$$$$$  /$$$$$$$
   | $$      | $$$$$   | $$  | $$ /$$__  $$ |____  $$| $$__  $$
   | $$      | $$__/   | $$  | $$| $$$$$$$$  /$$$$$$$| $$  \ $$
   | $$      | $$      | $$  | $$| $$_____/ /$$__  $$| $$  | $$
   | $$$$$$$$| $$$$$$$$| $$$$$$$/|  $$$$$$$|  $$$$$$$| $$  | $$
   |________/|________/|_______/  \_______/ \_______/|__/  |__/

   %s (%s) ver. %s
 
`, AUTOR, EMAIL, VERSION)
}
