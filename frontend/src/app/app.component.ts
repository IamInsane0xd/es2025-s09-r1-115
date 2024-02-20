import {Component, OnInit} from "@angular/core";
import { RouterOutlet } from '@angular/router';
import {StatsComponent} from "./stats/stats.component";
import {NgIf} from "@angular/common";
import {SearchComponent} from "./search/search.component";

export type ResponseContainer = {
  id: string;
  blockId: number;
  bayNum: number;
  stackNum: number;
  tierNum: number;
  arrivedAt: string;
}

export type DisplayedContainer = {
  id: string;
  blockId: number;
  bayNum: number;
  stackNum: number;
  tierNum: number;
  arrivedAt: string;
}

// noinspection ExceptionCaughtLocallyJS
@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, StatsComponent, NgIf, SearchComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.less'
})
export class AppComponent implements OnInit {
  title = 'frontend';
  currentPage = 0;
  public static containers: ResponseContainer[] = [];

  async ngOnInit() {
    try {
      const response = await fetch("http://localhost:3001/api/containers/search", {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`api server responded with status ${response.status}`);
      }

      AppComponent.containers = (await response.json()) as ResponseContainer[];

    } catch (e) {
      console.log(e);
    }
  }

  switchPage(page: number) {
    this.currentPage = page;
  }
}
