import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';

import { LedsService } from './leds/leds.service';
import { ButtonService } from './button/button.service';
import { UpdateService } from './update/update.service';

import { LedDisplayComponent } from './led-display/led-display.component';
import { ButtonComponent } from './button/button.component';
import { NavigationComponent } from './navigation/navigation.component';

@NgModule({
  declarations: [
    AppComponent,
    LedDisplayComponent,
    ButtonComponent,
    NavigationComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule
    // FormsModule,
  ],
  providers: [
    LedsService,
    ButtonService,
    UpdateService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
