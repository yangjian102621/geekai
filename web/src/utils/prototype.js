// add prototype method for array to insert item
Array.prototype.insert = function (index, item) {
  this.splice(index, 0, item)
}

Array.prototype.remove = function (index) {
  this.splice(index, 1)
}