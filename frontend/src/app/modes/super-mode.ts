import { deepEqual } from 'fast-equals';
import { deepCopy } from 'src/app/lib/deep-copy';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { CmdMode, Cmd } from 'src/app/websocket/commands';

export abstract class SuperMode {

  public backParameter: any
  public parameter: any
  public limits: any

  constructor(protected name: string, protected websocketService: WebsocketService) {
  }

  public getName() {
    return this.name
  }

  public getShortName() {
    return this.name.replace("Mode", "")
  }

  public receiveParameter(parm: any) {
    if (!deepEqual(this.backParameter, parm)) {
      this.backParameter = parm
      this.parameter = deepCopy(this.backParameter)
      setTimeout(() => M.updateTextFields());
    }
  }


  public receiveLimits(limits: any) {
    this.limits = limits
  }

  public sendParameter() {
    this.websocketService.send({
      cmd: "mode",
      parm: {
        id: this.name,
        parm: this.parameter
      } as CmdMode
    } as Cmd)
  }
}