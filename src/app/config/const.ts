export const REST_PROTOCOL = "http";
export const REST_ADRESS = window.location.hostname
export const REST_PORT = "2211";
export const REST_PREFIX_URL = REST_PROTOCOL + "://" + REST_ADRESS+":"+REST_PORT+"/"

export const REST_GET_SYSTEM_EXIT_URL = REST_PREFIX_URL + "exit"

export const REST_GET_LEDS_URL = REST_PREFIX_URL + "leds"
export const REST_GET_LEDS_COUNT_URL = REST_GET_LEDS_URL + "/count"
export const REST_GET_LEDS_ROWS_URL = REST_GET_LEDS_URL + "/rows"

export const REST_PRESS_SINGLE_URL = REST_PREFIX_URL + "press_single"
export const REST_PRESS_DOUBLE_URL = REST_PREFIX_URL + "press_double"
export const REST_PRESS_LONG_URL = REST_PREFIX_URL + "press_long"

export const REST_RANDOMIZE_URL = REST_PREFIX_URL + "mode/randomize"


export const REST_GET_MODE_URL = REST_PREFIX_URL + "mode"
export const REST_GET_MODE_RESOLVER_URL = REST_PREFIX_URL + "mode/resolver"

export const REST_MODE_SOLID_URL = REST_PREFIX_URL + "ModeSolid"
export const REST_MODE_SOLID_RAINBOW_URL = REST_PREFIX_URL + "ModeSolidRainbow"
export const REST_MODE_TRANSITION_RAINBOW_URL = REST_PREFIX_URL + "ModeTransitionRainbow"
export const REST_MODE_RUNNING_LED_URL = REST_PREFIX_URL + "ModeRunningLed"
export const REST_MODE_EMITTER_URL = REST_PREFIX_URL + "ModeEmitter"
