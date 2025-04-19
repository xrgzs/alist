package fs

import (
	"context"
	"path"

	"github.com/alist-org/alist/v3/internal/model"
)

type SyncOptions struct {
	LeftPath  string
	RightPath string
	Direction string   // "left_to_right", "right_to_left", "bidirectional"
	DiffTypes []string // "name", "size", "time"
	Recursive bool
	User      *model.User // 用于任务归属
	Ctx       context.Context
}

type SyncResult struct {
	Copied    []string
	Deleted   []string
	Skipped   []string
	Conflicts []string
}

// syncDir 递归同步单个目录
func syncDir(opt SyncOptions, relPath string, result *SyncResult) error {
	leftDir := path.Join(opt.LeftPath, relPath)
	rightDir := path.Join(opt.RightPath, relPath)
	leftObjs, _ := List(opt.Ctx, leftDir, &ListArgs{})
	rightObjs, _ := List(opt.Ctx, rightDir, &ListArgs{})
	leftMap := make(map[string]model.Obj)
	rightMap := make(map[string]model.Obj)
	for _, obj := range leftObjs {
		leftMap[obj.GetName()] = obj
	}
	for _, obj := range rightObjs {
		rightMap[obj.GetName()] = obj
	}
	// 处理左侧有右侧没有的
	for name, l := range leftMap {
		r, rok := rightMap[name]
		if !rok {
			if opt.Direction == "left_to_right" || opt.Direction == "bidirectional" {
				// 复制到右
				dstDir := rightDir
				_, err := Copy(opt.Ctx, path.Join(leftDir, name), dstDir)
				if err != nil {
					result.Conflicts = append(result.Conflicts, path.Join(relPath, name)+":"+err.Error())
				} else {
					result.Copied = append(result.Copied, path.Join(relPath, name))
				}
			}
			if l.IsDir() && opt.Recursive {
				syncDir(opt, path.Join(relPath, name), result)
			}
			continue
		}
		// 两边都有，判断是否需要同步
		if isDiffObj(l, r, opt.DiffTypes) {
			if opt.Direction == "left_to_right" {
				dstDir := rightDir
				_, err := Copy(opt.Ctx, path.Join(leftDir, name), dstDir)
				if err != nil {
					result.Conflicts = append(result.Conflicts, path.Join(relPath, name)+":"+err.Error())
				} else {
					result.Copied = append(result.Copied, path.Join(relPath, name))
				}
			} else if opt.Direction == "right_to_left" {
				dstDir := leftDir
				_, err := Copy(opt.Ctx, path.Join(rightDir, name), dstDir)
				if err != nil {
					result.Conflicts = append(result.Conflicts, path.Join(relPath, name)+":"+err.Error())
				} else {
					result.Copied = append(result.Copied, path.Join(relPath, name))
				}
			} else if opt.Direction == "bidirectional" {
				result.Conflicts = append(result.Conflicts, path.Join(relPath, name)+":双向差异需人工处理")
			}
		}
		if l.IsDir() && r.IsDir() && opt.Recursive {
			syncDir(opt, path.Join(relPath, name), result)
		}
	}
	// 处理右侧有左侧没有的
	for name, r := range rightMap {
		_, lok := leftMap[name]
		if !lok {
			if opt.Direction == "left_to_right" {
				// 删除右
				err := Remove(opt.Ctx, path.Join(rightDir, name))
				if err != nil {
					result.Conflicts = append(result.Conflicts, path.Join(relPath, name)+":"+err.Error())
				} else {
					result.Deleted = append(result.Deleted, path.Join(relPath, name))
				}
			} else if opt.Direction == "right_to_left" || opt.Direction == "bidirectional" {
				// 复制到左
				dstDir := leftDir
				_, err := Copy(opt.Ctx, path.Join(rightDir, name), dstDir)
				if err != nil {
					result.Conflicts = append(result.Conflicts, path.Join(relPath, name)+":"+err.Error())
				} else {
					result.Copied = append(result.Copied, path.Join(relPath, name))
				}
				if r.IsDir() && opt.Recursive {
					syncDir(opt, path.Join(relPath, name), result)
				}
			}
		}
	}
	return nil
}

// isDiffObj 判断两个对象是否有差异
func isDiffObj(l, r model.Obj, diffTypes []string) bool {
	for _, t := range diffTypes {
		switch t {
		case "name":
			if l.GetName() != r.GetName() {
				return true
			}
		case "size":
			if l.GetSize() != r.GetSize() {
				return true
			}
		case "time":
			if !l.ModTime().Equal(r.ModTime()) {
				return true
			}
		}
	}
	return false
}
