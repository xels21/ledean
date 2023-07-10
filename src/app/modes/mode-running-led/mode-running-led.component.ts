
import { Component, OnInit } from '@angular/core';
import { ModeRunningLedService } from './mode-running-led.service';




@Component({
  selector: 'app-mode-running-led',
  templateUrl: './mode-running-led.component.html',
  styleUrls: ['./mode-running-led.component.scss', '../../app.component.scss']
})
export class ModeRunningLedComponent implements OnInit {

  constructor(public service: ModeRunningLedService) {}

  ngOnInit(): void {
    // this.updateModeRunningLedParameter();
    // this.updateModeRunningLedLimits();
    // this.updateService.registerPolling({ cb: () => { this.updateModeRunningLedParameter() }, timeout: 500 })
    // setTimeout(() => { M.updateTextFields() }, 100);
  }


}
