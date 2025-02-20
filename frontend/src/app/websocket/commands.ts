import { RGB } from '../color/color';


export const CmdLedsId = "leds"
export const CmdLedsParameterId = "ledsParameter"
export const CmdButtonId = "button"
export const CmdModeId = "mode"
export const CmdModeLimitsId = "modeLimits"
export const CmdModeResolverId = "modeResolver"
export const CmdModeActionId = "action"
export const CmdModeActionRandomizeId = "randomize"
export const CmdModeActionPlayPause = "playPause"

export interface Cmd {
	cmd: string
	parm: any
}

export interface CmdLeds {
	leds: Array<RGB>
}

export interface CmdLedsParameter {
	rows: number
	count: number
}

export interface CmdButton {
	action: string
}

export interface CmdMode {
	id: string
	parm: any
}

export interface CmdModeLimits {
	id: string
	limits: any
}

export interface CmdModeResolver {
	modes: string[]
}



export interface CmdModeAction {
	action: string
}