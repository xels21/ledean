import { Injectable } from '@angular/core';
import { deepCopy } from 'src/app/lib/deep-copy';
import { deepEqual } from 'fast-equals';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { Cmd, CmdMode } from 'src/app/websocket/commands';
import { ParentMode } from 'src/app/modes/parent-mode';

export interface ModeTransitionRainbowParameter {
  roundTimeMs: number,
  brightness: number,
  spectrum: number,
  reverse: boolean,
}
export interface ModeTransitionRainbowLimits {
  minRoundTimeMs: number,
  maxRoundTimeMs: number,
  minBrightness: number,
  maxBrightness: number,
}

@Injectable({
  providedIn: 'root'
})
export class ModeTransitionRainbowService extends ParentMode {
  public backParameter: ModeTransitionRainbowParameter
  public parameter: ModeTransitionRainbowParameter
  public limits: ModeTransitionRainbowLimits

  constructor(protected websocketService: WebsocketService) {
    super("ModeTransitionRainbow", websocketService)
  }
}
