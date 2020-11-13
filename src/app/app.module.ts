import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import { LedDisplayComponent } from './led-display/led-display.component';
import { LedsService } from './leds/leds.service';
import { ButtonService } from './button/button.service';
import { ControlsComponent } from './controls/controls.component';

@NgModule({
  declarations: [
    AppComponent,
    LedDisplayComponent,
    ControlsComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule
    // FormsModule,
  ],
  providers: [
    LedsService,
    ButtonService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
