
import { Component, OnInit } from '@angular/core';
import { ModeTransitionRainbowService } from './mode-transition-rainbow.service';

@Component({
  selector: 'app-mode-transition-rainbow',
  templateUrl: './mode-transition-rainbow.component.html',
  styleUrls: ['./mode-transition-rainbow.component.scss','../../app.component.scss']
})
export class ModeTransitionRainbowComponent implements OnInit {


  constructor(public service:ModeTransitionRainbowService) {
  }

  ngOnInit(): void {
    // this.updateModeTransitionRainbowParameter();
    // this.updateModeTransitionRainbowLimits();
    // this.updateService.registerPolling({ cb: () => { this.updateModeTransitionRainbowParameter() }, timeout: 500 })
    // setTimeout( ()=>{M.updateTextFields()},100);
  }

}
