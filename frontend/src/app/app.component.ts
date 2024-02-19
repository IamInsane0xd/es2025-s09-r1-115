import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import {StatsComponent} from "./stats/stats.component";
import {NgIf} from "@angular/common";
import {SearchComponent} from "./search/search.component";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, StatsComponent, NgIf, SearchComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.less'
})
export class AppComponent {
  title = 'frontend';
  currentPage = 0;

  switchPage(page: number) {
    this.currentPage = page;
  }
}
