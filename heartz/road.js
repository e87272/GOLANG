var a = [1, 1, 1, 1, 1, 1, 1, 2, 2, 4, 4, 3, 3, 5, 2, 0, 1, 3, 4, 4, 2, 0, 1, 2, 1, 2, 1, 2]

function test() {
	console.log("123")
}

function road(rowNum, colNum, gameResult) {
	var resultArray = []
	var colArray = []
	var typeFlag = ""
	for (var i = 0; i < gameResult.length; i++) {
		if (typeFlag == "odd") {
			if (gameResult[i] % 2 == 1) {
				colArray.push(gameResult[i])
				typeFlag = "odd"
			} else {
				resultArray.push(colArray)
				colArray = []
				colArray.push(gameResult[i])
				typeFlag = "even"
			}
		} else if (typeFlag == "even") {
			if (gameResult[i] % 2 == 1) {
				resultArray.push(colArray)
				colArray = []
				colArray.push(gameResult[i])
				typeFlag = "odd"
			} else {
				colArray.push(gameResult[i])
				typeFlag = "even"
			}
		} else {
			if (gameResult[i] % 2 == 1) {
				colArray.push(gameResult[i])
				typeFlag = "odd"
			} else {
				colArray.push(gameResult[i])
				typeFlag = "even"
			}
		}
	}
	var roadArray = []
	for (var i = 0; i < resultArray.length; i++) {
		if (roadArray.length > 0) {
			if (roadArray[roadArray.length - 1][0] != "") {
				roadArray = roadPush(rowNum, roadArray, resultArray[i])
			} else {
				col = 1
				while (roadArray[roadArray.length - col][0] == "") {
					col++
				}
				col = col - 1
				row = 0
				while (roadArray[roadArray.length - col][row] == "") {
					row++
				}
				roadTemp = roadPush(row, [], resultArray[i])
				for (var j = 0; j < roadTemp.length; j++) {
					if (roadArray.length - col + j < roadArray.length) {
						for (var k = 0; k < roadTemp[j].length; k++) {
							roadArray[roadArray.length - col + j][k] = roadTemp[j][k]
						}
					} else {
						colArray = []
						for (var k = 0; k < rowNum; k++) {
							if (k < roadTemp[j].length) {
								colArray.push(roadTemp[j][k])
							} else {
								colArray.push("")
							}
						}
						roadArray.push(colArray)
					}
				}
			}
		} else {
			roadArray = roadPush(rowNum, roadArray, resultArray[i])
		}
	}

	//最後一行補空
	colArray = []
	for (var j = 0; j < rowNum; j++) {
		colArray.push("")
	}
	roadArray.push(colArray)

	return roadArray.slice(colNum * -1)
}

function roadPush(rowNum, roadArray, resultArray) {
	if (resultArray.length <= rowNum) {
		colArray = []
		for (var j = 0; j < rowNum; j++) {
			if (j < resultArray.length) {
				colArray.push(resultArray[j])
			} else {
				colArray.push("")
			}
		}
		roadArray.push(colArray)
	} else {
		roadArray.push(resultArray.slice(0, rowNum))
		for (var j = 0; j < resultArray.slice(rowNum).length; j++) {
			colArray = []
			for (var k = 0; k < rowNum - 1; k++) {
				colArray.push("")
			}
			colArray.push(resultArray[rowNum + j - 1])
			roadArray.push(colArray)
		}
	}

	return roadArray
}

test()