import {Component, OnInit} from "@angular/core";
import {NgForOf, NgIf} from "@angular/common";
import {YardComponent} from "./yard/yard.component";

type Stats = {
  blockId: number,
  capacity: number,
  averageAge: number,
  oldestContainerId: string,
  newestContainerId: string,
  emptyPositions: number,
  emptyBays: number,
  emptyStacks: number,
}

// noinspection ExceptionCaughtLocallyJS
@Component({
  selector: 'app-stats',
  standalone: true,
  imports: [
    NgIf,
    NgForOf,
    YardComponent,
  ],
  templateUrl: './stats.component.html',
  styleUrl: './stats.component.less'
})
export class StatsComponent implements OnInit {
  stats: Stats[] = []
  errorMsg = ""
  showErr = false

  async ngOnInit(): Promise<void> {
    this.showErr = false;

    try {
      const response = await fetch("http://localhost:3001/api/blocks/stat", {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`api server returned with status ${response.status}`)
      }

      this.stats = (await response.json()) as Stats[];

      if (this.stats == null) {
        throw new Error("api responded with null")
      }

    } catch (e) {
      this.showErr = true;

      if (e instanceof Error) {
        this.errorMsg = `error: ${e.message}`
        return
      }

      this.errorMsg = `error: ${e}`
      return
    }
  }
}
