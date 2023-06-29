import { RGB } from '../color/color';

export interface Cmd {
	cmd: string
	parm: any
}

export interface Cmd2cLeds {
	leds: Array<RGB>
}