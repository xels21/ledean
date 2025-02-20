import { Injectable } from '@angular/core';
import { deepCopy } from 'src/app/lib/deep-copy';
import { deepEqual } from 'fast-equals';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { Cmd, CmdMode } from 'src/app/websocket/commands';
import { SuperMode } from 'src/app/modes/super-mode';

export interface ModePictureParameter {
  // roundTimeMs: number,
  brightness: number,
  pictureChangeIntervallMs: number,
  pictureColumnNs: number,
  // spectrum: number,
  // reverse: boolean,
}
export interface ModePictureLimits {
  // minRoundTimeMs: number,
  // maxRoundTimeMs: number,
  minBrightness: number,
  maxBrightness: number,
  minPictureChangeIntervallMs: number,
  maxPictureChangeIntervallMs: number,
  minPictureColumnNs: number,
  maxPictureColumnNs: number,
}

@Injectable({
  providedIn: 'root'
})
export class ModePictureService extends SuperMode {
  public backParameter: ModePictureParameter
  public parameter: ModePictureParameter
  public limits: ModePictureLimits

  constructor(protected websocketService: WebsocketService) {
    super("ModePicture", websocketService)
  }
}
