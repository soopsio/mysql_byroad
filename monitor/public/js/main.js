var inspector = $.cookie('inspector');
if (!inspector) {
  $("#inspector option:first").attr('selected', true);
  $.cookie('inspector', $('#inspector').val(), { path: '/' });
} else {
  var flag = false;
  $('#inspector option').each(function () {
    if ($(this).val() == inspector) {
      $(this).attr('selected', true);
      flag = true;
      return false;
    }
  });
  if (!flag) {
    $("#inspector option:first").attr('selected', true);
    $.cookie('inspector', $('#inspector').val(), { path: '/' });
  }
}

function removeField(obj) {
  $(obj).closest('.form-group').parent().remove();
}

function removeFieldDiv(obj) {
  $(obj).parent().parent().remove();
}

function addRightField(schema, table, column) {
  var field = schema + '@@' + table + '@@' + column;
  var divid = schema + '-' + table;
  var data = '<div class="col-md-4" style="margin-bottom:5px;">\
  <div class="form-group">\
  <div class="input-group">\
  <input value="'+ column + '" readonly type="text" class="form-control">\
  <input name="fields" type="hidden" value="'+ field + '">\
  </div>\
  <div class="checkbox"><label><input type="checkbox" name="' + field + '" value="1"></input>是否推送值<label></div> \
  <a href="javascript:void(0)" onclick="removeField(this)"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></a>\
  </div>\
  </div>';
  addRightDiv(schema, table);
  if (!checkRightField(schema, table, column)) {
    $('#' + divid).find('.rowCont').append(data);
  }
}

function addRightDiv(schema, table) {
  var field = schema + '-' + table;
  var data = '<div class="col-md-12 right-field-div" id="' + field + '">\
  <h3><input class="form-control input-lg schema-name" onchange="changeFieldValue(this)" value="' + schema + '"></input> \
  <input class="form-control input-lg table-name" onchange="changeFieldValue(this)" value="' + table + '"></input>\
  <a href="javascript:void(0)" onclick="removeFieldDiv(this)" class="btn btn-lg">\
  <span class="glyphicon glyphicon-remove" aria-hidden="true"></span>\
  </a>\
  </h3>\
  <div class="row rowCont"></div></div>';
  var d = document.getElementById(field);
  if (!d) {
    $('#rightdiv').append(data);
  }
}

function checkRightField(schema, table, column) {
  var field = schema + '-' + table;
  var d = document.getElementById(field);
  return $(d).find(':text').is($(':text[value=' + column + ']'));
}

function changeFieldValue(obj) {
  var pardiv = $(obj).parents('.right-field-div')[0];
  var schema_val = $(pardiv).find('.schema-name').val();
  var table_val = $(pardiv).find('.table-name').val();
  var newVal = schema_val + '@@' + table_val;
  var hiddenobj = $(pardiv).find(':hidden');
  var checkboxs = $(pardiv).find(':checkbox')
  $(hiddenobj).each(function () {
    var fieldVal = $(this).prev().val();
    $(this).val(newVal + '@@' + fieldVal);
  });
  $(checkboxs).each(function () {
    var fieldVal = $(this).attr('name').split('@@')[2];
    $(this).attr('name', newVal + '@@' + fieldVal)
  });
}

function addTask() {
  if (confirm('确认添加？')) {
    var options = {
      url: '/task',
      type: 'post',
      dataType: 'json',
      data: $('#form').serialize(),
      success: function (data) {
        alert(data.Message);
      },
      error: function (data) {
        if (data.status == 422) {
          alert("任务数据格式错误");
        } else {
          alert(data.Message||'添加失败');
        }
      }
    };
    $.ajax(options);
    return false;
  }
}

function getTaskTopics() {
  var options = {
      url: '/task/getTopics',
      type: 'post',
      dataType: 'json',
      data: $('#form').serialize(),
      success: function (data) {
        alert('查找到的topic：' + data.Message);
      },
      error: function (data) {
        if (data.status == 422) {
          alert("任务数据格式错误");
        } else {
          alert(data.Message||'添加失败');
        }
      }
    };
    $.ajax(options);
}

function modifyTask() {
  if (confirm('确认修改？')) {
    var options = {
      url: '/task',
      type: 'put',
      dataType: 'json',
      data: $('#form').serialize(),
      success: function (data) {
        alert(data.Message);
      },
      error: function (data) {
        if (data.status == 422) {
          alert("任务数据格式错误");
        } else {
          alert(data.Message||'修改失败');
        }
      }
    };
    $.ajax(options);
    return false;
  }
}

function deleteTask(taskid) {
  if (confirm("确认删除？")) {
    var option = {
      type: 'delete',
      url: '/task/' + taskid,
      dataType: 'json',
      success: function (data) {
        location.reload();
      },
      error: function (data) {
        alert(data.Message||"操作失败");
      }
    };
    $.ajax(option);
  }
}

function changeTaskStat(taskid, stat) {
  var option = {
    type: 'post',
    url: 'task/changeStat/' + taskid,
    dataType: 'json',
    data: { "stat": stat },
    success: function (data) {
      location.reload();
    },
    error: function (data) {
      alert(data.Message||"操作失败");
    }
  };
  $.ajax(option);
}

var columnMap = {};

function getTables(schema) {
  $('#' + schema).find('.glyphicon').toggleClass('glyphicon-triangle-bottom glyphicon-triangle-right');
  if (columnMap[schema]) {
    $('#' + schema).next().collapse('toggle');
  } else {
    $.get('/tables', { schema: schema }, function (response) {
      columnMap[schema] = response;
      $('#' + schema).after(response);
      $('#' + schema).next().collapse('toggle');
    });
  }
}

function getColumns(schema, table) {
  var aid = '#' + schema + table;
  $(aid).find('.glyphicon').toggleClass('glyphicon-triangle-bottom glyphicon-triangle-right');
  if (columnMap[schema][table]) {
    $(aid).next().collapse('toggle');
  } else {
    $.get('/columns', { schema: schema, table: table }, function (response) {
      $(aid).after(response);
      $(aid).next().collapse('toggle');
      columnMap[schema] = {};
      columnMap[schema][table] = response;
    });
  }
}

function startSub(taskid) {
  var option = {
    type: 'post',
    url: 'task/'+taskid+'/startSub',
    dataType: 'json',
    success: function (data) {
      location.reload();
    },
    error: function (data) {
      alert(data.Message||"操作失败");
    }
  };
  $.ajax(option);
}

function stopSub(taskid) {
  var option = {
    type: 'post',
    url: 'task/'+taskid+'/stopSub',
    dataType: 'json',
    success: function (data) {
      location.reload();
    },
    error: function (data) {
      alert(data.Message||"操作失败");
    }
  };
  $.ajax(option);
}

function startPush(taskid) {
  var option = {
    type: 'post',
    url: 'task/'+taskid+'/startPush',
    dataType: 'json',
    success: function (data) {
      location.reload();
    },
    error: function (data) {
      alert(data.Message||"操作失败");
    }
  };
  $.ajax(option);
}

function stopPush(taskid) {
  var option = {
    type: 'post',
    url: 'task/'+taskid+'/stopPush',
    dataType: 'json',
    success: function (data) {
      location.reload();
    },
    error: function (data) {
      alert(data.Message||"操作失败");
    }
  };
  $.ajax(option);
}

$(function () {
  $('#pack-help').popover({
    trigger: 'hover',
    title: '数据封装协议',
    html: true,
    content: '<dd><dt>默认</dt><dl>旁路系统原有格式: 消息内容从post请求的body中读取。</dl><dl>消费方处理完成后返回 success</dl><dt>消息中心推送协议</dt><dl>使用消息中心的推送协议进行数据封装: message=POST["message"], jobid=GET["jobid"], retry_times=GET["retry_times"]</dl><dl>消费方处理完成后返回{"status": 1}</dl></dd>',
  });
  $("#regexp-help").popover({
    trigger: 'hover',
    title: '正则表达式',
    html: true,
    content: '<dd><dt>支持正则表达式</dt><dl>数据库名和表名都可以使用正则表达式</dl><dl>默认会在表达式前后添加<strong>"^"</strong>和<strong>"$"</strong>符号</dl><dl>之前的<strong>*</strong>的效果同现在的<strong>([\\w]+)</strong>一样</dl><dt>',
  });

  $('#search').jSearch({
    selector: '#column-list',
    child: 'li .schema-for-search',
    minValLength: 0,
    Found: function (elem) {
      $(elem).next().children().show();
    },
    NotFound: function (elem) {
      $(elem).parent().hide();
    },
    After: function (t) {
      if (!t.val().length) $('#column-list>li').show();
    }
  });

  $('.click-show').click(function () {
    $(this).toggleClass('click-show');
    $(this).parent().siblings().children().toggleClass('click-show');
  });

  $('#inspector').change(function () {
    $.cookie('inspector', $(this).val(), { path: '/' });
    location.reload();
  });
});


