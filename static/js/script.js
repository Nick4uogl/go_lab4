document.addEventListener("DOMContentLoaded", function () {
  // Перемикання вкладок
  const tabs = document.querySelectorAll(".tab");
  const calculators = document.querySelectorAll(".calculator");

  tabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      const tabId = tab.getAttribute("data-tab");

      // Оновлення активної вкладки
      tabs.forEach((t) => t.classList.remove("active"));
      tab.classList.add("active");

      // Показ активного калькулятора
      calculators.forEach((calc) => calc.classList.remove("active"));
      document.getElementById(tabId).classList.add("active");
    });
  });

  // Обробник для форми калькулятора кабелів
  document
    .getElementById("cable-form")
    .addEventListener("submit", function (e) {
      e.preventDefault();

      const formData = new FormData(this);

      fetch("/calculate/cable", {
        method: "POST",
        body: formData,
      })
        .then((response) => response.json())
        .then((data) => {
          document.getElementById("normal-current").textContent =
            data.normalCurrent + " A";
          document.getElementById("post-emergency-current").textContent =
            data.postEmergencyCurrent + " A";
          document.getElementById("economic-cross-section").textContent =
            data.economicCrossSection + " мм²";
          document.getElementById("minimum-cross-section").textContent =
            data.minimumCrossSection + " мм²";

          document.getElementById("cable-results").style.display = "block";
        })
        .catch((error) => {
          console.error("Error:", error);
        });
    });

  // Обробник для форми калькулятора короткого замикання
  document.getElementById("sc-form").addEventListener("submit", function (e) {
    e.preventDefault();

    const formData = new FormData(this);

    fetch("/calculate/shortcircuit", {
      method: "POST",
      body: formData,
    })
      .then((response) => response.json())
      .then((data) => {
        document.getElementById("reactor-impedance").textContent =
          data.reactorImpedance;
        document.getElementById("transformer-impedance").textContent =
          data.transformerImpedance;
        document.getElementById("total-impedance").textContent =
          data.totalImpedance;
        document.getElementById("initial-sc-current").textContent =
          data.initialShortCircuitCurrent + " A";

        document.getElementById("sc-results").style.display = "block";
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  });

  // Обробник для форми калькулятора мережі
  document
    .getElementById("network-form")
    .addEventListener("submit", function (e) {
      e.preventDefault();

      const formData = new FormData(this);

      fetch("/calculate/network", {
        method: "POST",
        body: formData,
      })
        .then((response) => response.json())
        .then((data) => {
          // 110kV шини
          document.getElementById("ish3").textContent = data.iSh3;
          document.getElementById("ish2").textContent = data.iSh2;
          document.getElementById("ish3-min").textContent = data.iSh3Min;
          document.getElementById("ish2-min").textContent = data.iSh2Min;

          // 10kV шини
          document.getElementById("ishn3").textContent = data.iShN3;
          document.getElementById("ishn2").textContent = data.iShN2;
          document.getElementById("ishn3-min").textContent = data.iShN3Min;
          document.getElementById("ishn2-min").textContent = data.iShN2Min;

          // Точка 10
          document.getElementById("iln3").textContent = data.iLN3;
          document.getElementById("iln2").textContent = data.iLN2;
          document.getElementById("iln3-min").textContent = data.iLN3Min;
          document.getElementById("iln2-min").textContent = data.iLN2Min;

          document.getElementById("network-results").style.display = "block";
        })
        .catch((error) => {
          console.error("Error:", error);
        });
    });
});
