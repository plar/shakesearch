const MAX_RESULTS_PER_QUERY = 20;

function escapeRegex(string) {
  return string.replace(/[/\-\\^$*+?.()|[\]{}]/g, '\\$&');
}

function highlightResults(q, rx, result) {
  let matches = [...result.matchAll(rx)];
  let indices = []
  matches.forEach((match) => {
      indices.push(match.index)
  });
  indices.reverse();

  for (let idx of indices) {
    result = result.substring(0, idx+q.length)+"</strong>"+result.substring(idx+q.length)
    result = result.substring(0, idx)+"<strong>"+result.substring(idx)
  }
  return result
}

const Controller = {
  search: (ev) => {
    ev.preventDefault();

    Controller.setOffset(0);

    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(data, results, false);
        Controller.setOffset((results && results.length) || 0);
      });
    });
  },

  setOffset: (newOffset) => {
    document.getElementById("offset").value = newOffset
  },

  loadMore: (ev) => {
    const data = Object.fromEntries(new FormData(form));
    fetch(`/search?q=${data.query}&offset=${data.offset}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(data, results, true);
        Controller.setOffset(parseInt(data.offset)+MAX_RESULTS_PER_QUERY);
      });
    });
  },

  updateTable: (data, results, append) => {
    const q = escapeRegex(data.query);
    const rx = new RegExp(`${q}`, 'gi')

    const rows = [];
    let num = parseInt(data.offset)
    for (let result of results) {
      num += 1
      result = highlightResults(q, rx, result)
      rows.push(`<tr><td>${num}.</td><td>...${result}...</td></tr>`);
    }

    const tableBody = rows.join("");
    const table = document.getElementById("table-body");
    if (append) {
      if (results && results.length) {
        table.innerHTML += tableBody;
      }
    } else {
      table.innerHTML = tableBody;
    }
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);

const loadMoreBtn = document.getElementById("load-more");
loadMoreBtn.addEventListener("click", Controller.loadMore);