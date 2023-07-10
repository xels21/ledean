import { Component, OnInit } from '@angular/core';
import { ModesService } from './modes.service'
import { HttpClient } from '@angular/common/http';
import { REST_RANDOMIZE_URL } from '../config/const';
import { Cmd, CmdModeAction, CmdModeActionId, CmdModeActionRandomizeId } from '../websocket/commands';
import { WebsocketService } from '../websocket/websocket.service';

@Component({
  selector: 'app-modes',
  templateUrl: './modes.component.html',
  styleUrls: ['./modes.component.scss']
})
export class ModesComponent implements OnInit {

  constructor(public modesService: ModesService, private websocketService: WebsocketService, private httpClient: HttpClient) { }

  ngOnInit(): void {
    // $('ul.tabs').tabs();
    var instance = M.Tabs.init($('ul.tabs'), { /*swipeable: true*/ })[0];
    $('.indicator').css('background-color', 'teal');
    // setTimeout(()=>instance.select(this.modesService.modeResolver[this.modesService.activeMode]),100)
    this.modesService.setOnModeChange(() => {
      instance.select(this.modesService.modeResolver[this.modesService.activeMode])
    })
  }

  randomizeMode() {
    this.websocketService.send({
      cmd: CmdModeActionId,
      parm: { action: CmdModeActionRandomizeId } as CmdModeAction
    } as Cmd)

    // this.httpClient.get(REST_RANDOMIZE_URL).subscribe();
  }

}
