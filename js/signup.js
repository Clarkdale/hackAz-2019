(function() {

	window.onload = function() {

		// Events for each tile 
		var donator = document.getElementById("donatorOrCharity1");
		var charity = document.getElementById("donatorOrCharity2");

		donator.onclick = moreForms;
		charity.onclick = moreForms;
	};

	function moreForms() {
		var charity = document.getElementById("donatorOrCharity2");
		var form = document.getElementById("charity-donator");
		if(charity.checked) {

			var button = document.getElementById("signup-btn");
    		button.parentNode.removeChild(button);

			var div = document.createElement("div");
			div.className = "form-group";
			div.id = "website-form";

			var input = document.createElement("input");
			input.type = "text";
			input.className = "form-control";
			input.id = "website-link";
			input.placeholder = "Enter website link";

			div.appendChild(input);
			form.appendChild(div);

			div = document.createElement("div");
			div.className = "form-group";
			div.id = "charity-form";

			input = document.createElement("input");
			input.type = "text";
			input.className = "form-control";
			input.id = "charity-name";
			input.placeholder = "Enter charity name";

			div.appendChild(input);
			form.appendChild(div);

			var a = document.createElement("a");
			a.className = "btn btn-primary btn-xl";
			a.id = "signup-btn";
			a.href = "#"; // TODO: add actual link
			a.text = "Sign Up";

			form.appendChild(a);
		} else {

			var inp1 = document.getElementById("website-link");
    		inp1.parentNode.removeChild(inp1);

    		var inp2 = document.getElementById("charity-name");
    		inp2.parentNode.removeChild(inp2);
		}
	}

})();