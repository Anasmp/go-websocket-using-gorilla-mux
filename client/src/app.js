import { Todo } from "./todo";

export class App {
  constructor() {
   
    this.heading = 'Todos';
    this.todos = [];
    this.todoDescription = '';
    this.username =''
    this.loggedIn = false

  }

  addTodo() {
    if (this.todoDescription) {
      var todo = new Todo(this.todoDescription);
      var msg = {"type":"add","todo":todo,"username":this.username}
      this.socket.send(JSON.stringify(msg));
      this.todoDescription = '';
    }
  }
doLogin() {
  var app = this
  if(!this.username){
    return 
  }
  this.loggedIn = true
  this.socket =new WebSocket('ws://localhost:8081/ws');

  this.socket.addEventListener('error',function(event){
    console.log(event)
  })
  this.socket.addEventListener('open',function(event){
    var msg = {"type":"hello","username":app.username}
    app.socket.send(JSON.stringify(msg))
  })
  this.socket.addEventListener('message',function(event) {
    var msg = JSON.parse(event.data)
    app.todos = msg.todos
  })
}
removeTodo(id) {
  // console.log(todo)
    var msg = {"type":"delete","id":id,"username":this.username}
    this.socket.send(JSON.stringify(msg))
  }

  toggleTodoDone(id) {
    // console.log("updateTodo()",todo)
    var msg = {"type":"toggle.done","id":id,"username":this.username}
    this.socket.send(JSON.stringify(msg))
  }



}

