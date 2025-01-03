package goes

import "github.com/olivere/elastic/v7"

// 索引 - 所有索引名称
func (cli *GoEs) IndexNames() ([]string, error) {
	return cli.cli.IndexNames()
}

// 索引 - 是否存在
func (cli *GoEs) IndexExists(index string) (bool, error) {
	return cli.cli.IndexExists(index).Do(cli.ctx)
}

// 索引 - 创建
func (cli *GoEs) IndexCreate(index, body string) (err error) {
	var exist bool
	exist, err = cli.IndexExists(index)
	if err != nil {
		return
	}
	if exist {
		err = nil
		return
	}
	_, err = cli.cli.CreateIndex(index).BodyString(body).Do(cli.ctx)
	return
}

// 索引 - 查询 - 文档结构、索引设置、data
func (cli *GoEs) IndexGet(index string) (*elastic.IndicesGetResponse, error) {
	resp, err := cli.cli.IndexGet().Index(index).Do(cli.ctx)
	if err != nil {
		return nil, err
	}
	return resp[index], nil
}

// 索引 - 查看 - 文档结构
func (cli *GoEs) IndexMapping(index string) (map[string]interface{}, error) {
	return cli.cli.GetMapping().Index(index).Do(cli.ctx)
}

// 索引 - 修改 - 文档结构
func (cli *GoEs) IndexUpdateMapping(index, body string) (*elastic.PutMappingResponse, error) {
	return cli.cli.PutMapping().Index(index).BodyString(body).Do(cli.ctx)
}

// 索引 - 查看 - 索引设置
func (cli *GoEs) IndexSettings(index string) (map[string]*elastic.IndicesGetSettingsResponse, error) {
	return cli.cli.IndexGetSettings().Index(index).Do(cli.ctx)
}

// 索引 - 修改 - 索引设置
func (cli *GoEs) IndexUpdateSettings(index, body string) (*elastic.IndicesPutSettingsResponse, error) {
	return cli.cli.IndexPutSettings().Index(index).BodyString(body).Do(cli.ctx)
}

// 索引 - 别名 - 添加
func (cli *GoEs) IndexAlias(index, aliasName string) (*elastic.AliasResult, error) {
	return cli.cli.Alias().Add(index, aliasName).Do(cli.ctx)
}

// 索引 - 别名 - 删除
func (cli *GoEs) IndexAliasRemove(index, aliasName string) (*elastic.AliasResult, error) {
	return cli.cli.Alias().Remove(index, aliasName).Do(cli.ctx)
}

// 索引 - 删除
func (cli *GoEs) IndexDelete(index string) (err error) {
	var exist bool
	exist, err = cli.IndexExists(index)
	if err != nil {
		return
	}
	if !exist {
		err = nil
		return
	}
	_, err = cli.cli.DeleteIndex(index).Do(cli.ctx)
	return
}
