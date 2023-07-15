import { Injectable } from '@angular/core';
import { WebsocketService } from '../websocket/websocket.service';
import { Cmd, CmdButton,CmdButtonId } from '../websocket/commands';


@Injectable({
  providedIn: 'root'
})
export class ButtonService {

  constructor(private websocketService: WebsocketService) { }

  public pressSingle() {
    this.websocketService.send({
      cmd: CmdButtonId,
      parm: { action: "single" } as CmdButton
    } as Cmd)
  }

  public pressDouble() {
        this.websocketService.send({
      cmd: CmdButtonId,
      parm: { action: "double" } as CmdButton
    } as Cmd)
  }

  public pressLong() {
        this.websocketService.send({
      cmd: CmdButtonId,
      parm: { action: "long" } as CmdButton
    } as Cmd)
  }
}
