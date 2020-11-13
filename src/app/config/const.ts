export const REST_PROTOCOL = "http";
export const REST_ADRESS = window.location.hostname
export const REST_PORT = "2211";
export const REST_PREFIX_URL = REST_PROTOCOL + "://" + REST_ADRESS+":"+REST_PORT+"/"

export const REST_GET_LEDS_URL = REST_PREFIX_URL + "leds"

export const REST_PRESS_SINGLE_URL = REST_PREFIX_URL + "press_single"
export const REST_PRESS_DOUBLE_URL = REST_PREFIX_URL + "press_double"
export const REST_PRESS_LONG_URL = REST_PREFIX_URL + "press_long"