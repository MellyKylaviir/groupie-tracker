const item = document.querySelectorAll(".item");

for (i = 0; i < item.length; i++) {
    item[i].addEventListener('mouseover', (e) => {
        console.log(e);
        e.currentTarget.style.borderWidth = "5px";
        e.currentTarget.style.boxShadow = "10px 10px 20px darkblue";
    });
}

for (i = 0; i < item.length; i++) {
    item[i].addEventListener('mouseout', (e) => {
        console.log(e);
        e.currentTarget.style.borderWidth = "";
        e.currentTarget.style.boxShadow = "";
    });
}