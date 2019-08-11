function setAsideMenu(json) {
	const ul = document.getElementById("actions-menu")

	json.actions.map(action => {
		const a = document.createElement("a")
		a.href = `${document.location.href}#${action.request}`
		a.innerText = action.name

		const li = document.createElement("li")

		li.appendChild(a)
		ul.appendChild(li)
	})

	const footer = document.querySelector("footer")
	footer.innerHTML = `Source code: <a href="${json.project_url}">github.com</a>`
}

function setHeader(json) {
	const header = document.querySelector("header")
	header.id = "introduction"

	const h1 = document.createElement("h1")
	h1.innerText = json.title

	const description = document.createElement("div")
	description.innerHTML = json.description

	header.appendChild(h1)
	header.appendChild(description)
}

function setActions(json) {
	const main = document.querySelector("main")

	const actions = document.createElement("div")
	actions.id = "actions"

	actions.classList.add("actions")
	const actionTitle = document.createElement("h3")
	actionTitle.innerText = "Actions"

	actions.appendChild(actionTitle)

	json.actions.map(action => {
		const title = document.createElement("h4")
		title.innerText = action.name
		title.id = action.request

		const description = document.createElement("p")
		description.innerText = action.description

		const requestTitle = document.createElement("h5")
		requestTitle.innerText = "Request"

		const example = document.createElement("h6")
		example.classList.add("example")
		example.innerText = "Example"

		const request = json.requests[action.request]

		const pre = document.createElement("pre")
		pre.innerText = JSON.stringify(request.example, undefined, 4)

		const detailsDiv = request.details.map(detail => {
			const div = document.createElement("div")
			div.classList.add("details")

			const detailTitle = document.createElement("h6")
			detailTitle.innerText = detail.key

			if (detail.mandatory) {
				detailTitle.innerText += "*"
			}

			const detailDescription = document.createElement("p")
			detailDescription.innerText = detail.description

			div.appendChild(detailTitle)
			div.appendChild(detailDescription)
			return div
		})

		const responseTitle = document.createElement("h5")
		responseTitle.innerText = "Response"

		const successfulDot = document.createElement("div")
		successfulDot.classList.add("successful-dot")
		const successfulText = document.createElement("h6")
		successfulText.innerText = "Successful"

		const successful = document.createElement("div")
		successful.classList.add("line")
		successful.appendChild(successfulDot)
		successful.appendChild(successfulText)

		const responsesOK = json.responses[action.response_ok]
		const responses = responsesOK.map(res => {
			const resDiv = document.createElement("div")
			resDiv.classList.add("successful-response")

			if (res.detail.target_title) {
				const successfulTitle = document.createElement("h5")
				successfulTitle.innerText = res.detail.target_title
				resDiv.appendChild(successfulTitle)
			}

			const detailSuccessfulRes = document.createElement("p")
			detailSuccessfulRes.innerText = res.detail.reason

			const exampleSuccessful = document.createElement("h6")
			exampleSuccessful.classList.add("example")
			exampleSuccessful.style.marginLeft = "25px"
			exampleSuccessful.innerText = "Example"

			const pre_successful = document.createElement("pre")
			pre_successful.innerText = JSON.stringify(res.example, undefined, 4)

			const detailSuccessfulResTarget = document.createElement("h6")
			detailSuccessfulResTarget.classList.add("margin-top")
			detailSuccessfulResTarget.innerText = "TARGET"

			const detailSuccessfulResTargetDescription = document.createElement("p")
			detailSuccessfulResTargetDescription.classList.add("margin-top")
			detailSuccessfulResTargetDescription.innerText = res.detail.target

			resDiv.appendChild(exampleSuccessful)
			resDiv.appendChild(pre_successful)
			resDiv.appendChild(pre_successful)
			resDiv.appendChild(detailSuccessfulResTarget)
			resDiv.appendChild(detailSuccessfulResTargetDescription)

			return resDiv
		})

		const errorDot = document.createElement("div")
		errorDot.classList.add("error-dot")
		const errorText = document.createElement("h5")
		errorText.innerText = "Error"

		const error = document.createElement("div")
		error.classList.add("line")
		error.appendChild(errorDot)
		error.appendChild(errorText)

		const errors = action.response_error.map(err => {
			const errJSON = json.responses[err]

			const divErr = document.createElement("div")
			divErr.classList.add("err-json")

			const divErrBtn = document.createElement("div")
			const divErrBtnText = document.createElement("span")
			divErrBtnText.innerText = errJSON.detail.reason

			const divErrBtnIcon = document.createElement("span")
			divErrBtnIcon.innerText = ">"
			divErrBtnIcon.classList.add("error-icon")

			divErrBtn.appendChild(divErrBtnIcon)
			divErrBtn.appendChild(divErrBtnText)

			const pre_err = document.createElement("pre")
			pre_err.innerText = JSON.stringify(errJSON.example, undefined, 4)

			const detailErrResTarget = document.createElement("h6")
			detailErrResTarget.classList.add("margin-top")
			detailErrResTarget.innerText = "TARGET"

			const detailErrResTargetDescription = document.createElement("p")
			detailErrResTargetDescription.innerText = errJSON.detail.target

			divErr.classList.add("hide")

			divErrBtn.onclick = () => {
				if (divErr.classList.contains("hide")) {
					divErr.classList.remove("hide")
					divErr.querySelector(".error-icon").style.transform = "rotate(90deg)"
				} else {
					divErr.classList.add("hide")
					divErr.querySelector(".error-icon").style.transform = "rotate(0deg)"
				}
			}

			divErr.appendChild(divErrBtn)
			divErr.appendChild(pre_err)
			divErr.appendChild(detailErrResTarget)
			divErr.appendChild(detailErrResTargetDescription)
			return divErr
		})

		actions.appendChild(title)
		actions.appendChild(description)
		actions.appendChild(requestTitle)
		actions.appendChild(example)
		actions.appendChild(pre)
		detailsDiv.map(d => {
			actions.appendChild(d)
		})
		actions.appendChild(responseTitle)
		actions.appendChild(successful)
		responses.map(res => {
			actions.appendChild(res)
		})
		actions.appendChild(error)
		errors.map(e => {
			actions.appendChild(e)
		})
	})

	main.appendChild(actions)
}

const JSON_URL = "https://api.myjson.com/bins/mbe3n"

fetch(JSON_URL).then(res => {
	res.json().then(json => {
		document.title = json.title
		setAsideMenu(json)
		setHeader(json)
		setActions(json)
	})
})
