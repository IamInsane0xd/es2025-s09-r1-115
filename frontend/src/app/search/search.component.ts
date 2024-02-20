import {Component, OnInit} from "@angular/core";
import {NgForOf, NgIf} from "@angular/common";
import {ResponseContainer, DisplayedContainer} from "../app.component";

// noinspection ExceptionCaughtLocallyJS
@Component({
  selector: 'app-search',
  standalone: true,
  imports: [
    NgIf,
    NgForOf,
  ],
  templateUrl: './search.component.html',
  styleUrl: './search.component.less'
})
export class SearchComponent implements OnInit {
  searchResult: DisplayedContainer[] = [];
  errorMsg = "";
  showErr = false;

  ngOnInit() {
    (document.querySelector("#search-form") as HTMLFormElement).addEventListener("submit",
      (e) => {
      e.preventDefault();
    });

    (document.querySelectorAll(".form-input") as NodeListOf<HTMLInputElement>).forEach((e) => {
      e.addEventListener("change", () => {
        this.search().then(r => r);
      });
    });

    this.search().then(r => r);
  }

  async search() {
    this.searchResult = [];
    this.showErr = false;

    let reqUrl = "http://localhost:3001/api/containers/search";
    let first = true;
    const id = (document.querySelector("#form-id") as HTMLInputElement).value;
    const blockId = Number((document.querySelector("#form-block-id") as HTMLInputElement).value);
    const bayNum = Number((document.querySelector("#form-bay-num") as HTMLInputElement).value);
    const stackNum = Number((document.querySelector("#form-stack-num") as HTMLInputElement).value);
    const tierNum = Number((document.querySelector("#form-tier-num") as HTMLInputElement).value);
    const sortBy = (document.querySelector("#form-sort-by") as HTMLSelectElement).value;

    if (id.trim().length != 0) {
      reqUrl += `${first ? '?' : '&'}id=${id}`;
      first = false;
    }

    if (blockId !== 0) {
      reqUrl += `${first ? '?' : '&'}blockId=${blockId}`;
      first = false;
    }

    if (bayNum !== 0) {
      reqUrl += `${first ? '?' : '&'}bayNum=${bayNum}`;
      first = false;
    }

    if (stackNum !== 0) {
      reqUrl += `${first ? '?' : '&'}stackNum=${stackNum}`;
      first = false;
    }

    if (tierNum !== 0) {
      reqUrl += `${first ? '?' : '&'}tierNum=${tierNum}`;
      first = false;
    }

    if (sortBy !== "none") {
      reqUrl += `${first ? '?' : '&'}sortBy=${sortBy}`;
    }

    console.log(reqUrl);

    try {
      const response = await fetch(reqUrl, {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
      });

      if (!response.ok) {
        throw new Error("no containers found")
      }

      const containers = (await response.json()) as ResponseContainer[];

      for (let container of containers) {
        let unix = false;
        let unixNum = 0;

        if (container.arrivedAt.length === 13) {
          unix = true;
          unixNum = Number(container.arrivedAt);
        }

        this.searchResult.push({
          id: container.id,
          blockId: container.blockId,
          bayNum: container.bayNum,
          stackNum: container.stackNum,
          tierNum: container.tierNum,
          arrivedAt: (unix ? new Date(unixNum) : new Date(container.arrivedAt)).toLocaleString(),
        });
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
