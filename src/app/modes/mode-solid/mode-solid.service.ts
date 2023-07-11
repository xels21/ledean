import { Injectable } from '@angular/core';
import { RGB } from 'src/app/color/color'
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { ParentMode } from '../parent-mode';


export interface ModeSolidParameter {
  rgb: RGB,
  brightness: number,
}
export interface ModeSolidLimits {
  minBrightness: number,
  maxBrightness: number,
}

@Injectable({
  providedIn: 'root'
})
export class ModeSolidService extends ParentMode {
  public backParameter: ModeSolidParameter
  public parameter: ModeSolidParameter
  public limits: ModeSolidLimits

  constructor(protected websocketService: WebsocketService) {
    super("ModeSolid", websocketService)
  }
}
