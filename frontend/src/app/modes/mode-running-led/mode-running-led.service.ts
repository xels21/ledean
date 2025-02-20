import { Injectable } from '@angular/core';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { SuperMode } from 'src/app/modes/super-mode';



// type RunningLedStyle = "linear" | "trigonometric"
export enum RunningLedStyle {
  LINEAR = "linear",
  TRIGONOMETRIC = "trigonometric",
}

export interface ModeRunningLedParameter {
  brightness: number,
  fadePct: number,
  roundTimeMs: number,
  hueFrom: number,
  huerTo: number,
  style: RunningLedStyle,
}

export interface ModeRunningLedLimits {
  minRoundTimeMs: number,
  maxRoundTimeMs: number,
  minBrightness: number,
  maxBrightness: number,
}


@Injectable({
  providedIn: 'root'
})
export class ModeRunningLedService extends SuperMode{
  public backParameter: ModeRunningLedParameter
  public parameter: ModeRunningLedParameter
  public limits: ModeRunningLedLimits

  constructor(protected websocketService: WebsocketService) { 
    super("ModeRunningLed",websocketService)
  }

  getAllStyles() {
    return new Array<RunningLedStyle>(
      RunningLedStyle.LINEAR,
      RunningLedStyle.TRIGONOMETRIC,
    )
  }
}
