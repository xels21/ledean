import { Component, OnInit } from '@angular/core';
import { ModesService } from './modes.service'

import * as M from 'materialize-css';
import * as $ from "jquery";
// import $ from "jquery";

@Component({
  selector: 'app-modes',
  templateUrl: './modes.component.html',
  styleUrls: ['./modes.component.scss']
})
export class ModesComponent implements OnInit {

  constructor(public modesService: ModesService) { }

  ngOnInit(): void {
    // $('ul.tabs').tabs();
    var instance = M.Tabs.init($('ul.tabs'), { swipeable: true });
    $('.indicator').css('background-color', 'teal');
    this.modesService.setOnModeChange(() => { 
      // instance.updateTabIndicator();
    })
    // Or with jQuery

    // $(document).ready(function(){
    //   $('.tabs').tabs();
    // });
  }



}
