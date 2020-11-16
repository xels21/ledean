import { Component,  } from '@angular/core';
import { LedsService } from '../leds/leds.service';


@Component({
  selector: 'app-led-display',
  templateUrl: './led-display.component.html',
  styleUrls: ['./led-display.component.scss']
})
export class LedDisplayComponent {

  constructor(public ledsService:LedsService) { }

  // onPollingChanged(e){
  //   this.ledsService.pollingActive = e.target.checked
  //   this.ledsService.checkPollingInterval()
  // }

}
