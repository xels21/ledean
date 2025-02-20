import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { WebsocketService } from '../websocket/websocket.service';
import { ModeEmitterService } from '../modes/mode-emitter/mode-emitter.service';
import { ModeGradientService } from '../modes/mode-gradient/mode-gradient.service';
import { ModeSpectrumService } from '../modes/mode-spectrum/mode-spectrum.service';
import { ModeRunningLedService } from '../modes/mode-running-led/mode-running-led.service';
import { ModeSolidService } from '../modes/mode-solid/mode-solid.service';
import { ModeSolidRainbowService } from '../modes/mode-solid-rainbow/mode-solid-rainbow.service';
import { ModePictureService } from '../modes/mode-picture/mode-picture.service';
import { ModeTransitionRainbowService } from '../modes/mode-transition-rainbow/mode-transition-rainbow.service';
import { Cmd, CmdMode, CmdModeId, CmdModeLimits, CmdModeResolver } from '../websocket/commands';
import { SuperMode } from './super-mode';


@Injectable({
  providedIn: 'root'
})
export class ModesService {

  activeMode: number
  isPictureMode: boolean
  public modeResolver: Array<string>
  onModeChange: () => any
  //modeServices: Array<SuperMode>

  constructor(private httpClient: HttpClient, private websocketService: WebsocketService
    , public modeEmitterService: ModeEmitterService
    , public modeGradientService: ModeGradientService
    , public modeSpectrumService: ModeSpectrumService
    , public modeRunningLedService: ModeRunningLedService
    , public modeSolidRainbowService: ModeSolidRainbowService
    , public modeSolidService: ModeSolidService
    , public modePictureService: ModePictureService
    , public modeTransitionRainbowService: ModeTransitionRainbowService
  ) {

    // this.modeServices = [modeEmitterService, modeGradientService, modeRunningLedService, modeSolidRainbowService, modeSolidService, modeTransitionRainbowService]
    this.setOnModeChange(() => { })
    this.websocketService.modeChanged.subscribe((cmdMode: CmdMode) => this.modeChanged(cmdMode))
    this.websocketService.modeLimitChanged.subscribe((cmdModeLimits: CmdModeLimits) => this.modeLimitChanged(cmdModeLimits))
    this.websocketService.modeResolverChanged.subscribe((cmdModeResolver: CmdModeResolver) => { this.modeResolver = cmdModeResolver.modes; })
  }

  public isModeAvailable(mode: string) : boolean{
    if(!this.modeResolver){return false}
    return this.modeResolver.indexOf(mode) != -1
  } 

  private modeChanged(cmdMode: CmdMode) {
    console.log(cmdMode)
    switch (cmdMode.id) {
      case this.modeEmitterService.getName():
        this.modeEmitterService.receiveParameter(cmdMode.parm)
        break;
      case this.modeGradientService.getName():
        this.modeGradientService.receiveParameter(cmdMode.parm)
        break;
      case this.modeSpectrumService.getName():
        this.modeSpectrumService.receiveParameter(cmdMode.parm)
        break;
      case this.modeRunningLedService.getName():
        this.modeRunningLedService.receiveParameter(cmdMode.parm)
        break;
      case this.modeSolidRainbowService.getName():
        this.modeSolidRainbowService.receiveParameter(cmdMode.parm)
        break;
      case this.modeSolidService.getName():
        this.modeSolidService.receiveParameter(cmdMode.parm)
        break;
      case this.modeTransitionRainbowService.getName():
        this.modeTransitionRainbowService.receiveParameter(cmdMode.parm)
        break;
        case this.modePictureService.getName():
          this.modePictureService.receiveParameter(cmdMode.parm)
          break;
      default:
        console.log("Unknown limit change for mode: ", cmdMode.id)
        break;
    }
    this.updateActiveModeString(cmdMode.id)
  }

  private modeLimitChanged(cmdModeLimits: CmdModeLimits) {
    switch (cmdModeLimits.id) {
      case this.modeEmitterService.getName():
        this.modeEmitterService.receiveLimits(cmdModeLimits.limits)
        break;
      case this.modeGradientService.getName():
        this.modeGradientService.receiveLimits(cmdModeLimits.limits)
        break;
      case this.modeSpectrumService.getName():
        this.modeSpectrumService.receiveLimits(cmdModeLimits.limits)
        break;
      case this.modeRunningLedService.getName():
        this.modeRunningLedService.receiveLimits(cmdModeLimits.limits)
        break;
      case this.modeSolidService.getName():
        this.modeSolidService.receiveLimits(cmdModeLimits.limits)
        break;
      case this.modeSolidRainbowService.getName():
        this.modeSolidRainbowService.receiveLimits(cmdModeLimits.limits)
        break;
      case this.modeTransitionRainbowService.getName():
        this.modeTransitionRainbowService.receiveLimits(cmdModeLimits.limits)
        break;
        case this.modePictureService.getName():
          this.modePictureService.receiveLimits(cmdModeLimits.limits)
          break;

      default:
        console.log("Unknown limit change for mode: ", cmdModeLimits.id)
        break;
    }
  }


  public setOnModeChange(onModeChange: () => any) {
    this.onModeChange = onModeChange
  }

  public updateActiveModeString(name: string) {
    var id = this.modeStrToIdx(name)
    if (id == undefined) {
      console.log("mode resolver not ready yet: ", name)
      setTimeout(() => this.updateActiveModeString(name), 100)
    } else {
      this.updateActiveMode(id)
    }
  }

  public updateActiveMode(id: number) {
    if (this.activeMode != id) {
      this.activeMode = id
      if (this.activeMode >= 0) {
        this.isPictureMode = false
        this.onModeChange()
      } else {
        this.isPictureMode = true
      }
    }
  }

  public isActive(mode: string) {
    return mode == this.modeResolver[this.activeMode]
  }

  public switchState(mode: string) {
    this.websocketService.send({
      cmd: CmdModeId,
      parm: {
        id: mode,
        parm: undefined
      } as CmdMode
    } as Cmd)
  }

  public modeStrToIdx(mode: string) {
    if (this.modeResolver == undefined) {
      return undefined
    }
    return this.modeResolver.findIndex(m => { return m == mode })
  }



}
