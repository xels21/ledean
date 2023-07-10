
import { Component, OnInit } from '@angular/core';
import { ModeSolidRainbowService } from './mode-solid-rainbow.service';




@Component({
  selector: 'app-mode-solid-rainbow',
  templateUrl: './mode-solid-rainbow.component.html',
  styleUrls: ['./mode-solid-rainbow.component.scss', '../../app.component.scss']
})
export class ModeSolidRainbowComponent implements OnInit {


  constructor(public service: ModeSolidRainbowService) {}

  ngOnInit(): void {
    // this.updateModeSolidRainbowParameter();
    // this.updateModeSolidRainbowLimits();
    // this.updateService.registerPolling({ cb: () => { this.updateModeSolidRainbowParameter() }, timeout: 500 })
    // setTimeout(() => { M.updateTextFields() }, 100);
  }

}
