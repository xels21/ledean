import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { REST_PRESS_SINGLE_URL, REST_PRESS_DOUBLE_URL, REST_PRESS_LONG_URL } from '../config/const';

@Injectable({
  providedIn: 'root'
})
export class ButtonService {

  constructor(private httpClient: HttpClient) { }

  public pressSingle() {
    this.httpClient.get(REST_PRESS_SINGLE_URL).subscribe();
  }

  public pressDouble() {
    this.httpClient.get(REST_PRESS_DOUBLE_URL).subscribe();
  }

  public pressLong() {
    this.httpClient.get(REST_PRESS_LONG_URL).subscribe();
  }
}
