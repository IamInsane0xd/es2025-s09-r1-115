import { Component } from '@angular/core';
import {NgForOf, NgIf} from "@angular/common";
import {AppComponent} from "../../app.component";

@Component({
  selector: 'app-yard',
  standalone: true,
  imports: [
    NgIf,
    NgForOf,
  ],
  templateUrl: './yard.component.html',
  styleUrl: './yard.component.less'
})
export class YardComponent {
  currentBlock = 1;
  blockLoop = [1, 2, 3, 4];
  bayLoop = [1, 2, 3, 4, 5];
  stackLoop = [1, 2, 3, 4, 5];
  tierLoop = [1, 2, 3, 4, 5];

  switchBlock(block: number) {
    this.currentBlock = block;
  }

  checkContainer(blockId: number, bayNum: number, stackNum: number, tierNum: number): boolean {
    for (let container of AppComponent.containers) {
      if (container.blockId == blockId && container.bayNum == bayNum &&
          container.stackNum == stackNum && container.tierNum == tierNum) {
        return true;
      }
    }

    return false;
  }

  protected readonly AppComponent = AppComponent;
}
