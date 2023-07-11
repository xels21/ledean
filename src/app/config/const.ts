export const REST_PROTOCOL = "http";
export const REST_ADRESS = window.location.hostname
export const REST_PORT = "2211";
export const REST_PREFIX_URL = REST_PROTOCOL + "://" + REST_ADRESS+":"+REST_PORT+"/"

export const REST_GET_SYSTEM_EXIT_URL = REST_PREFIX_URL + "exit"