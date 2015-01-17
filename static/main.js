
$(function(){
	GoSSH.getData(GoSSH.populatePage);
});


var GoSSH = {
	initialGet: true,
	getData: function(method){
		 $.ajax({url:'/data',dataType:'json',success:function(result){
    		method(result);
  		}});		
	},
	populatePage: function(data){
		if (data == null)
			return;
		$.each(data, function(k, v) {
    		$($("#tmpl").html()).appendTo(document.body)
			$("div.gridster:last").attr("data",GoSSH.tidyClass(k)).find("div.server").text(k).end().find("ul").addClass(GoSSH.tidyClass(k));
			createGrid("ul."+GoSSH.tidyClass(k),v);
  		});
		window.setInterval("GoSSH.getData(GoSSH.refreshPage)", 5000);
	},
	refreshPage:function(){
		GoSSH.getData(GoSSH.updatePanels);
	},
	updatePanels: function(data){
		$.each(data, function(k, v) {
			var klass = GoSSH.tidyClass(k);
			var parent = $("div.gridster[data="+klass+"]");
			
			$.each(v.Values, function(vk, vv) {
				$(parent).find("li."+vk.toLowerCase()).find("span:first").html(vk + "<br/>" + vv);
			});
		});
	},
	tidyClass: function(original){
		return original.split(".").join("");
	}
};

function createGrid(gridName, data){
	var gridster = $(".gridster " + gridName).gridster({
		helper: 'clone',
		widget_margins: [5, 5],
		widget_base_dimensions: [140, 140],
		min_rows: 3,max_cols:2,
		max_size_x:2,
		resize:{enabled:true}
	}).data('gridster');
	$.each(data.Values, function(k, v) {
		var key = k.toLowerCase();
		gridster.add_widget("<li class='"+key+"'><span class='t'>"+k+"<br/>"+v+"</span></li>");
	});
	/*
	gridster.add_widget('<li class="ram new"><span class="t">Free RAM<br/>866MB</span></li>');
	gridster.add_widget('<li class="load new"><span class="t">Load Average<br/>0.01</span></li>',2,1);
	gridster.add_widget('<li class="network new"><span class="t">Network<br/>11Mbps</span></li>');
	gridster.add_widget('<li class="disk new"><span class="t">Disk Usage<br/>18%</span></li>');
	*/
}