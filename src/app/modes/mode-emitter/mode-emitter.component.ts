import { Component, OnInit } from '@angular/core';

import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';

import { ModeEmitterService } from './mode-emitter.service';



@Component({
  selector: 'app-mode-emitter',
  templateUrl: './mode-emitter.component.html',
  styleUrls: ['./mode-emitter.component.scss','../../app.component.scss']
})
export class ModeEmitterComponent implements OnInit {


  constructor(private updateService: UpdateService, private httpClient: HttpClient, public service:ModeEmitterService) { }

  ngOnInit(): void {
    // this.updateModeEmitterParameter();
    // this.updateModeEmitterLimits();
    // this.updateService.registerPolling({ cb: () => { this.updateModeEmitterParameter() }, timeout: 500 })
  }


}
