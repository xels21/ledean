import { Component, OnInit } from '@angular/core';
import { ModesService } from './modes.service'
import { HttpClient } from '@angular/common/http';
import { REST_RANDOMIZE_URL } from '../config/const';

import * as M from 'materialize-css';
import * as $ from "jquery";
// import $ from "jquery";

@Component({
  selector: 'app-modes',
  templateUrl: './modes.component.html',
  styleUrls: ['./modes.component.scss']
})
export class ModesComponent implements OnInit {

  constructor(public modesService: ModesService, private httpClient:HttpClient) { }

  ngOnInit(): void {
    // $('ul.tabs').tabs();
    var instance = M.Tabs.init($('ul.tabs'), { /*swipeable: true*/ })[0];
    $('.indicator').css('background-color', 'teal');
    // setTimeout(()=>instance.select(this.modesService.modeResolver[this.modesService.activeMode]),100)
    this.modesService.setOnModeChange(()=>{
      instance.select(this.modesService.modeResolver[this.modesService.activeMode])
    })
  }

  randomizeMode(){
    this.httpClient.get(REST_RANDOMIZE_URL).subscribe();

}

}
