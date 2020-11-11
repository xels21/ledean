import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import { LedDisplayComponent } from './led-display/led-display.component';
import { LedsService } from './leds/leds.service';

@NgModule({
  declarations: [
    AppComponent,
    LedDisplayComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule
    // FormsModule,
  ],
  providers: [
    LedsService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
