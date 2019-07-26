const searchBtn = document.querySelector("button[type='submit']");
const input = document.querySelector("input");

input.addEventListener("keyup", filterCourts);
// searchBtn.addEventListener("click", searchCourt);

function filterCourts() {
    const filteredCourts = courts.filter(court => {
        const regexp =  new RegExp(input.value, "gi");
        return court.nom.match(regexp) || court.arrondissement.match(regexp);
    });
    container.innerHTML = "";
    if(filteredCourts.length === 0){
      notFound();
      return;
    }
    filteredCourts.map(court => createCard(court));
  }
function notFound(){
    const div = document.createElement("div");
    div.classList.add("no-result")
    div.innerHTML = `
    <p>Aucun terrain trouvé</p>
        `;
    container.appendChild(div);
  }

function createCard(court) {
    const card = document.createElement("a");
    card.href = `/courts/${court.ID}`;
    card.innerHTML = `
        <div class="card-body">
            <h4 class="card-title">${court.nom}</h4>
            <ul>
                <li>
                    Dimensions: ${court.dimensions} m2
                </li>
                <li>
                    Revetement: ${court.revetement}
                </li>
                <li>
                    Decouvert: ${court.decouvert}
                </li>
                <li>
                    Eclairage: ${court.eclairage}
                </li>
            </ul>
        </div>
        `;
    container.appendChild(card);
  }