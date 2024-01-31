function CreateCluster(username, name, password, database, replicas) {
  let xhr = new XMLHttpRequest();
  xhr.open("POST", "/api/create");
  xhr.setRequestHeader("Accept", "application/json");
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.withCredentials = true;
  
  xhr.onreadystatechange = function () {
    if (xhr.readyState === 4) {
      console.log(xhr.status);
      console.log(xhr.responseText);
    }};
  
  let data = `{
    "username": "${username}",
    "name": "${name}",
    "password": "${password}",
    "replicas": ${replicas},
    "database": "${database}"
  }`;
  
  xhr.send(data);
}

function DeleteCluster(name) {
  let xhr = new XMLHttpRequest();
  xhr.open("POST", "/api/delete");
  xhr.setRequestHeader("Accept", "application/json");
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.withCredentials = true;
  
  xhr.onreadystatechange = function () {
    if (xhr.readyState === 4) {
      console.log(xhr.status);
      console.log(xhr.responseText);
    }};
  
  let data = `{
    "name": "${name}"
  }`;
  
  xhr.send(data);
}

function UpdateCluster(name, replicas) {
  let xhr = new XMLHttpRequest();
  xhr.open("POST", "/api/update");
  xhr.setRequestHeader("Accept", "application/json");
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.withCredentials = true;
  
  xhr.onreadystatechange = function () {
    if (xhr.readyState === 4) {
      console.log(xhr.status);
      console.log(xhr.responseText);
    }};
  
  let data = `{
    "name": "${name}",
    "replicas": ${replicas}
  }`;
  
  xhr.send(data);
}
