import { Component, OnInit } from '@angular/core';
import { LedsService } from '../leds/leds.service';

@Component({
  selector: 'app-led-display',
  templateUrl: './led-display.component.html',
  styleUrls: ['./led-display.component.scss']
})
export class LedDisplayComponent implements OnInit {


  constructor(public ledsService:LedsService) { }

  ngOnInit(): void {
  }

  onPollingChanged(e){
    this.ledsService.pollingActive = e.target.checked
    this.ledsService.checkPollingInterval()
  }

}
